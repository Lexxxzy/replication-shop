import random
from functools import lru_cache

from locust import (
    FastHttpUser,
    SequentialTaskSet,
    task,
    run_single_user,
    between,
)
from locust.contrib.fasthttp import FastResponse
from mimesis import Person, Address

import models


class SequenceOfTasks(SequentialTaskSet):
    current_user: models.User | None = None
    cart: models.MyCart | None = None
    products: list[models.Product] | None = None
    products_by_name: list[models.Product] | None = None
    orders: models.MyOrder | None = None

    headers = {"Content-Type": "application/json"}

    person = Person("en")
    address = Address()

    def on_start(self) -> None:
        self.register()
        self.wait()
        self.login()

    def on_stop(self):
        self.logout()

    def register(self):
        current_user = models.User(
            name=self.person.full_name(),
            email=self.person.email(unique=True),
            password=self.person.password(length=6),
        )
        with self.client.post(
                "/register",
                json=current_user.model_dump(include={"name", "email", "password"}),
        ) as resp:
            resp: FastResponse
            resp.raise_for_status()
            self.current_user = current_user

    def login(self):
        payload = {
            "email": self.current_user.email,
            "password": self.current_user.password,
        }
        with self.client.post(
                "/login",
                json=payload,
        ) as resp:
            resp: FastResponse
            resp.raise_for_status()
            self.headers.update({"session": resp.headers.get("Set-Cookie")})

    def logout(self):
        self.client.post("/logout", headers=self.headers)

    @task
    def get_products(self):
        with self.client.get("/products") as resp:
            resp: FastResponse
            resp.raise_for_status()
            data = resp.json()
            self.products = [models.Product(**x) for x in data.get("products")]

    def get_products_func(self, title: str):
        with self.client.get(f"/products?title={title}") as resp:
            resp: FastResponse
            resp.raise_for_status()
            data = resp.json()
            self.products_by_name = [models.Product(**x) for x in data.get("products")]

    @lru_cache
    def products_list(self):
        return [
            self.products[0],
            self.products[len(self.products) // 2],
            self.products[len(self.products) - 1],
        ]

    @task
    def get_product_by_name(self):
        if self.products:
            for product in self.products_list():
                self.wait()
                self.get_products_func(product.name)

    @task
    def get_cart(self):
        with self.client.get(
                "/my/cart", headers=self.headers
        ) as resp:
            resp: FastResponse
            resp.raise_for_status()
            self.cart = models.MyCart(**resp.json())

    @task
    def add_to_cart(self):
        product: models.Product = random.choice(self.products)
        payload = {
            "item_id": product.id,
            "quantity": random.randint(1, 10),
        }
        self.client.put("/my/cart/add", json=payload, headers=self.headers)
        self.get_cart()

    @task
    def remove_from_cart(self):
        if self.cart is not None:
            if self.cart.cart is not None:
                payload = {"item_id": self.cart.cart[-1].id}
                self.client.delete(
                    "/my/cart/remove", json=payload, headers=self.headers
                )
                self.get_cart()

    @task
    def get_orders(self):
        with self.client.get("/my/orders", headers=self.headers) as resp:
            resp: FastResponse
            resp.raise_for_status()
            data = resp.json()
            self.orders = models.MyOrder(**data)

    @task
    def add_order(self):
        self.add_to_cart()
        self.wait()
        payload = {"delivery_address": self.address.address()}
        self.client.post("/my/orders/add", json=payload, headers=self.headers)
        self.get_orders()


class ApiUser(FastHttpUser):
    host = "http://localhost:80"
    wait_time = between(1, 2)
    tasks = [SequenceOfTasks]


if __name__ == "__main__":
    run_single_user(ApiUser)
