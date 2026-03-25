namespace ProxyPatternLab.Models;

public class Book
{
    public int Id { get; set; }
    public string Title { get; set; } = string.Empty;
    public string Author { get; set; } = string.Empty;
    public string Genre { get; set; } = string.Empty;
    public int Year { get; set; }
    public double Rating { get; set; }
}

public class ServiceLog
{
    public string Timestamp { get; set; } = string.Empty;
    public string Message { get; set; } = string.Empty;
    public string Type { get; set; } = "info"; // info | cache | warn | db
}

public class ServiceResponse
{
    public List<Book> Books { get; set; } = new();
    public List<ServiceLog> Logs { get; set; } = new();
    public long ElapsedMs { get; set; }
    public bool FromCache { get; set; }
}
