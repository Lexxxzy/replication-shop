namespace ShopClient;

public class Logger(string logFilePath)
{
    string LogFilePath { get; set; } = logFilePath;

    public void LogToFile(string type, string endpoint, string method, string responseStatus, long responseDurationMs, string readOrWrite)
    {
        var logEntry = $"{DateTime.Now:yyyy-MM-dd HH:mm:ss.fffff},{type},{endpoint},{method},{responseStatus},{responseDurationMs},{readOrWrite}";

        using var file = new StreamWriter(LogFilePath, true);
        file.WriteLine(logEntry);
    }
}