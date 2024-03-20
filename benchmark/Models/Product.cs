using Newtonsoft.Json;

namespace ShopClient.Models;

public class Product
{
    [JsonProperty("id")]
    public int Id { get; set; }

    [JsonProperty("name")]
    public string Name { get; set; }

    [JsonProperty("price")]
    public float Price { get; set; }

    [JsonProperty("manufacturer")]
    public string Manufacturer { get; set; }

    [JsonProperty("type_name")]
    public string TypeName { get; set; }
}

