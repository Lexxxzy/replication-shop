using Newtonsoft.Json;

namespace ShopClient.Models;

public class NginxInstance
{
    [JsonProperty("ip")] public string Ip { get; set; }
    [JsonProperty("port")] public int Port { get; set; }

    public string GetHttpsAddress() => $"http://{Ip}:{Port}";

    public HttpClient GetHttpClient()
    {
        var client = new HttpClient();
        client.BaseAddress = new Uri(GetHttpsAddress());
        return client;
    }
}