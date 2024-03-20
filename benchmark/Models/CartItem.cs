using Newtonsoft.Json;

namespace ShopClient.Models;

public class ShoppingCart
{
    [JsonProperty("cart")]
    public List<CartItem> Cart { get; set; }
    
    [JsonProperty("total")]
    public float Total { get; set; }
}

public class CartItem
{
    [JsonProperty("id")]
    public int Id { get; set; }
    
    [JsonProperty("product")]
    public string Product { get; set; }
    
    [JsonProperty("price")]
    public float Price { get; set; }
    
    [JsonProperty("quantity")]
    public int Quantity { get; set; }
}