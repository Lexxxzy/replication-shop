using System.Diagnostics;
using System.Net.Http.Json;
using System.Text;
using Bogus;
using Newtonsoft.Json;

namespace ShopClient.Models;

public class User
{
    [JsonProperty("name")] public string Name { get; set; } 

    [JsonProperty("email")] public string Email { get; set; }

    [JsonProperty("password")] public string Password { get; set; }

    [JsonProperty("delivery_address")] public string DeliveryAddress { get; set; }

    [JsonProperty("token")] public string? Token { get; set; }

    private Logger? Logger { get; set; }

    public User()
    {
        var faker = new Faker();
        Name = faker.Name.FullName();
        Email = faker.Internet.Email();
        Password = faker.Internet.Password();
        DeliveryAddress = faker.Address.FullAddress();
    }

    public void SetLogger(Logger logger)
    {
        Logger = logger;
    }

    private async Task<HttpResponseMessage> SendRequestWithLoggingAsync(HttpClient client, HttpRequestMessage request)
    {
        if (!string.IsNullOrEmpty(Token))
        {
            request.Headers.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue("session", Token);
        }

        Logger.LogToFile($"Request: {request.Method} {request.RequestUri}");
        var stopwatch = Stopwatch.StartNew();
        var response = await client.SendAsync(request);
        stopwatch.Stop();

        Logger.LogToFile(
            $"Response: {(int)response.StatusCode} {response.ReasonPhrase} (Duration: {stopwatch.ElapsedMilliseconds} ms)");
        return response;
    }

    public async Task<HttpResponseMessage> RegisterUserAsync(HttpClient client)
    {
        var request = new HttpRequestMessage(HttpMethod.Post, "register")
        {
            Content = JsonContent.Create(new { Name, Email, Password })
        };
        return await SendRequestWithLoggingAsync(client, request);
    }

    public async Task<HttpResponseMessage> LoginUserAsync(HttpClient client)
    {
        var request = new HttpRequestMessage(HttpMethod.Post, "login")
        {
            Content = JsonContent.Create(new { Email, Password })
        };
        var response = await SendRequestWithLoggingAsync(client, request);
        Token = response.Headers.GetValues("Set-Cookie").First();
        return response;
    }

    public async Task<IEnumerable<Product>> GetAllProductsAsync(HttpClient client)
    {
        var request = new HttpRequestMessage(HttpMethod.Get, "products");
        var response = await SendRequestWithLoggingAsync(client, request);
        var productsRaw = await response.Content.ReadAsStringAsync();
        var products = JsonConvert.DeserializeObject<Dictionary<string, List<Product>>>(productsRaw)!["products"];
        return products;
    }

    public async Task<Product> GetProductAsync(HttpClient client, string productName)
    {
        var request = new HttpRequestMessage(HttpMethod.Get, $"products?title={productName}");
        var response = await SendRequestWithLoggingAsync(client, request);
        var productRaw = await response.Content.ReadAsStringAsync();
        var product = JsonConvert.DeserializeObject<Dictionary<string, List<Product>>>(productRaw)!["products"][0];
        return product;
    }

    public async Task<HttpResponseMessage> AddProductToCartAsync(HttpClient client, int productId)
    {
        var request = new HttpRequestMessage(HttpMethod.Put, "my/cart/add")
        {
            Content = new StringContent(JsonConvert.SerializeObject(new { item_id = productId, quantity = 1 }),
                Encoding.UTF8, "application/json")
        };
        return await SendRequestWithLoggingAsync(client, request);
    }

    public async Task<HttpResponseMessage> RemoveProductFromCartAsync(HttpClient client, int productId)
    {
        var request = new HttpRequestMessage(HttpMethod.Delete, "my/cart/remove")
        {
            Content = new StringContent(JsonConvert.SerializeObject(new { item_id = productId }), Encoding.UTF8,
                "application/json")
        };
        return await SendRequestWithLoggingAsync(client, request);
    }

    public async Task<List<CartItem>> GetCartAsync(HttpClient client)
    {
        var request = new HttpRequestMessage(HttpMethod.Get, "my/cart");
        var response = await SendRequestWithLoggingAsync(client, request);
        var responseContent = await response.Content.ReadAsStringAsync();

        var cart = JsonConvert.DeserializeObject<ShoppingCart>(responseContent) ?? new ShoppingCart();
        return cart.Cart;
    }

    public async Task<HttpResponseMessage> GetOrderAsync(HttpClient client)
    {
        var request = new HttpRequestMessage(HttpMethod.Get, "my/orders");
        return await SendRequestWithLoggingAsync(client, request);
    }

    public async Task<int> AddOrderAsync(HttpClient client, string deliveryAddress)
    {
        var request = new HttpRequestMessage(HttpMethod.Post, "my/orders/add")
        {
            Content = new FormUrlEncodedContent(new[]
            {
                new KeyValuePair<string, string>("delivery_address", deliveryAddress)
            })
        };
        var response = await SendRequestWithLoggingAsync(client, request);

        var responseContent = await response.Content.ReadAsStringAsync();
        var orderId = JsonConvert.DeserializeObject<Dictionary<string, int>>(responseContent)!["order_id"];
        return orderId;
    }

    public async Task<HttpResponseMessage> CancelOrderAsync(HttpClient client, int orderId)
    {
        var request = new HttpRequestMessage(HttpMethod.Delete, "my/orders/cancel")
        {
            Content = new StringContent(JsonConvert.SerializeObject(new { order_id = orderId }), Encoding.UTF8,
                "application/json")
        };
        return await SendRequestWithLoggingAsync(client, request);
    }

    public async Task<HttpResponseMessage> LogoutUserAsync(HttpClient client)
    {
        var request = new HttpRequestMessage(HttpMethod.Post, "logout");
        return await SendRequestWithLoggingAsync(client, request);
    }
}