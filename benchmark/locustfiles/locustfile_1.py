import random

from locust import FastHttpUser, SequentialTaskSet, task, between, run_single_user
from locust.contrib.fasthttp import FastResponse
from mimesis import Person, Text

import models

person = Person("en")
text = Text()


class ApiUser(FastHttpUser):
    host = "http://localhost:80"
    wait_time = between(1, 3)

    @task
    class SequenceOfTasks(SequentialTaskSet):
        current_user: models.User | None = None
        cart: models.MyCart | None = None
        products: list[models.Product] | None = None
        orders: models.MyOrder | None = None
        headers = {"Content-Type": "application/json"}

        @task
        def register(self):
            current_user = models.User(
                name=person.full_name(),
                email=person.email(unique=True),
                password=person.password(length=6),
            )
            with self.client.post(
                "/register",
                json=current_user.model_dump(include={"name", "email", "password"}),
            ) as resp:
                resp: FastResponse
                resp.raise_for_status()
                self.current_user = current_user

        @task
        def login(self):
            if not self.current_user:
                self.wait()
                self.register()
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

        @task
        def get_products(self):
            with self.client.get("/products") as resp:
                resp: FastResponse
                resp.raise_for_status()
                data = resp.json()
                self.products = [models.Product(**x) for x in data.get("products")]

        @task
        def get_cart(self):
            with self.client.get("/my/cart", headers=self.headers) as resp:
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
            self.get_products()
            self.add_to_cart()
            payload = {"delivery_address": text.text()}
            self.client.post("/my/orders/add", json=payload, headers=self.headers)
            self.get_orders()

        @task
        def cancel_order(self):
            order: models.Order = random.choice(self.orders.orders)
            payload = {"order_id": order.id}
            self.client.delete("/my/orders/cancel", json=payload, headers=self.headers)
            self.get_orders()

        @task
        def logout(self):
            self.client.post("/logout", headers=self.headers)
            self.cart = None
            self.orders = None
            self.products = None
            self.current_user = None


if __name__ == "__main__":
    run_single_user(ApiUser)
