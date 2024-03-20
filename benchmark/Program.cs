using System.Text.Json;
using CommandLine;
using ShopClient;
using ShopClient.Models;

var options = Parser.Default.ParseArguments<Options>(args);

var instances = await ReadNginxConfig(options.Value.NginxConfigPath);
if (instances == null)
{
    Console.WriteLine("Failed to read nginx config file, or it was empty. Exiting...");
    return;
}

var clients = new Queue<HttpClient>(instances.Select(instance => instance.GetHttpClient()));
var rnd = new Random();
var logger = new Logger(options.Value.LogfilePath);

while (true)
{
    try
    {
        // Round robin for client
        var client = clients.Dequeue();
        Console.WriteLine("Using client: " + client.BaseAddress);

        await SpawnUserAndMakeRequests(client, rnd);

        clients.Enqueue(client);
    }
    catch (OperationCanceledException)
    {
        Console.WriteLine("Operation canceled");
        break;
    }
    catch (Exception ex)
    {
        Console.WriteLine(ex.Message);
        break;
    }
}

async Task SpawnUserAndMakeRequests(HttpClient httpClient, Random random)
{
    {
        var user = new User();
        user.SetLogger(logger);

        await user.RegisterUserAsync(httpClient);
        await user.LoginUserAsync(httpClient);

        var products = (await user.GetAllProductsAsync(httpClient)).ToList();
        // NOTE: Add random number of products to cart (from 4 to 100 -> can add same products multiple times)
        for (int i = 0; i < random.Next(4, 50); i++)
        {
            var product = await user.GetProductAsync(httpClient, products[random.Next(0, products.Count)].Name);
            await user.AddProductToCartAsync(httpClient, product.Id);
        }

        var cart = await user.GetCartAsync(httpClient);

        // NOTE: Remove random number of products from cart
        // TASK: Предположение: модель пользователя-шопоголика подразумевает “муки выбора” и частое изменение корзины.
        var randomNumberOfIdsToDelete = GetRandomNumberOfIds(cart);
        foreach (var id in randomNumberOfIdsToDelete)
        {
            await user.RemoveProductFromCartAsync(httpClient, id);
        }

        // NOTE: Place order
        var orderId = await user.AddOrderAsync(httpClient, user.DeliveryAddress);
        await user.GetOrderAsync(httpClient);

        // NOTE: Cancel order with 50% chance
        if (random.NextDouble() > 0.5)
        {
            await user.CancelOrderAsync(httpClient, orderId);
        }

        // NOTE: Logout user to clear session from backend
        await user.LogoutUserAsync(httpClient);
    }
    
    static List<int> GetRandomNumberOfIds(List<CartItem> cartItems)
    {
        var rnd = new Random();
        int listLength = rnd.Next(0, cartItems.Count - 1); // Generate random length for the list

        List<int> idsList = new List<int>();
        foreach (var item in cartItems)
        {
            for (int i = 0; i < item.Quantity; i++)
            {
                idsList.Add(item.Id); // Add item ID based on its quantity
            }
        }

        // Shuffle the list to randomize which IDs are included
        idsList = idsList.OrderBy(x => rnd.Next()).ToList();

        // Trim the list to the randomly determined length
        if (idsList.Count > listLength)
        {
            idsList = idsList.Take(listLength).ToList();
        }

        return idsList;
    }
}

static async Task<List<NginxInstance>?> ReadNginxConfig(string path)
{
    try
    {
        string jsonContent = await File.ReadAllTextAsync(path);
        var options = new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true
        };

        var instances = JsonSerializer.Deserialize<List<NginxInstance>>(jsonContent, options)!;

        Console.WriteLine("Using nginx instances:");
        foreach (var instance in instances)
        {
            Console.WriteLine($"- {instance.GetHttpsAddress()}");
        }

        return instances;
    }
    catch (IOException e)
    {
        Console.WriteLine($"An error occurred while reading the file: {e.Message}");
    }
    catch (JsonException e)
    {
        Console.WriteLine($"An error occurred while deserializing the file content: {e.Message}");
    }

    return null;
}