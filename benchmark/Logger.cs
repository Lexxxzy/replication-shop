namespace ShopClient;

public class Logger(string logFilePath)
{
    string LogFilePath { get; set; } = logFilePath;

    public void LogToFile(string message)
    {
        var logEntry = $"{DateTime.Now:yyyy-MM-dd HH:mm:ss.fffff}: {message}\n";
        File.AppendAllText(LogFilePath, logEntry);
    }
}