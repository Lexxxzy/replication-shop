namespace ShopClient.Models;

using CommandLine;

public class Options
{
    [Option('c', "config", Required = true, HelpText = "Path to the nginx config file.")]
    public string NginxConfigPath { get; set; }

    [Option('l', "log", Required = true, HelpText = "Path to the log file.")]
    public string LogfilePath { get; set; }
}
