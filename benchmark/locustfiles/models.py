from pydantic import BaseModel


class Product(BaseModel):
    id: int | None = None
    name: str
    price: float
    manufacturer: str | None = None
    type_name: str | None = None


class ProductInMyCart(BaseModel):
    id: int | None = None
    product: str
    price: float
    quantity: int


class MyCart(BaseModel):
    cart: list[ProductInMyCart] | None = None
    total: float


class Order(BaseModel):
    id: int | None = None
    delivery_address: str | None = None
    order_date: str | None = None
    total_price: float
    cart_items: list[ProductInMyCart] | None = None


class MyOrder(BaseModel):
    orders: list[Order] | None = None


class User(BaseModel):
    id: str | None = None
    name: str
    email: str
    password: str
    created_at: str | None = None
