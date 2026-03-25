using System.Diagnostics;
using ProxyPatternLab.Models;

namespace ProxyPatternLab.Services;



public class RealBookService : IBookService
{
    private readonly List<ServiceLog> _logs = new();

    private static readonly List<Book> _database = new()
    {
        new Book { Id = 1, Title = "Clean Code",              Author = "Robert C. Martin", Genre = "Programming", Year = 2008, Rating = 4.7 },
        new Book { Id = 2, Title = "Design Patterns",         Author = "Gang of Four",     Genre = "Programming", Year = 1994, Rating = 4.8 },
        new Book { Id = 3, Title = "The Pragmatic Programmer",Author = "Hunt & Thomas",    Genre = "Programming", Year = 1999, Rating = 4.6 },
        new Book { Id = 4, Title = "Dune",                    Author = "Frank Herbert",    Genre = "Sci-Fi",      Year = 1965, Rating = 4.9 },
        new Book { Id = 5, Title = "Foundation",              Author = "Isaac Asimov",     Genre = "Sci-Fi",      Year = 1951, Rating = 4.7 },
        new Book { Id = 6, Title = "Neuromancer",             Author = "William Gibson",   Genre = "Sci-Fi",      Year = 1984, Rating = 4.5 },
        new Book { Id = 7, Title = "1984",                    Author = "George Orwell",    Genre = "Dystopia",    Year = 1949, Rating = 4.8 },
        new Book { Id = 8, Title = "Brave New World",         Author = "Aldous Huxley",   Genre = "Dystopia",    Year = 1932, Rating = 4.6 },
    };

    private void Log(string message, string type = "db")
    {
        _logs.Add(new ServiceLog
        {
            Timestamp = DateTime.Now.ToString("HH:mm:ss.fff"),
            Message   = message,
            Type      = type
        });
    }

    public async Task<ServiceResponse> GetAllBooksAsync()
    {
        var sw = Stopwatch.StartNew();
        _logs.Clear();

        Log("Opening database connection...");
        await Task.Delay(600); // моделирование задержки базы данных

        Log("Executing: SELECT * FROM Books");
        await Task.Delay(800);

        Log($"Query returned {_database.Count} rows", "info");
        Log("Closing database connection...");

        sw.Stop();
        return new ServiceResponse { Books = _database.ToList(), Logs = _logs.ToList(), ElapsedMs = sw.ElapsedMilliseconds, FromCache = false };
    }

    public async Task<ServiceResponse> GetBookByIdAsync(int id)
    {
        var sw = Stopwatch.StartNew();
        _logs.Clear();

        Log("Opening database connection...");
        await Task.Delay(500);

        Log($"Executing: SELECT * FROM Books WHERE Id = {id}");
        await Task.Delay(600);

        var book = _database.FirstOrDefault(b => b.Id == id);
        Log(book is null ? "No rows found" : "Query returned 1 row", "info");
        Log("Closing database connection...");

        sw.Stop();
        var books = book is null ? new List<Book>() : new List<Book> { book };
        return new ServiceResponse { Books = books, Logs = _logs.ToList(), ElapsedMs = sw.ElapsedMilliseconds, FromCache = false };
    }

    public async Task<ServiceResponse> GetBooksByGenreAsync(string genre)
    {
        var sw = Stopwatch.StartNew();
        _logs.Clear();

        Log("Opening database connection...");
        await Task.Delay(500);

        Log($"Executing: SELECT * FROM Books WHERE Genre = '{genre}'");
        await Task.Delay(700);

        var result = _database.Where(b => b.Genre.Equals(genre, StringComparison.OrdinalIgnoreCase)).ToList();
        Log($"Query returned {result.Count} rows", "info");
        Log("Closing database connection...");

        sw.Stop();
        return new ServiceResponse { Books = result, Logs = _logs.ToList(), ElapsedMs = sw.ElapsedMilliseconds, FromCache = false };
    }

    public void ResetLogs() => _logs.Clear();
}
