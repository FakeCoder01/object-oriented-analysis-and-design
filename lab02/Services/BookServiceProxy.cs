using System.Diagnostics;
using ProxyPatternLab.Models;

namespace ProxyPatternLab.Services;


public class BookServiceProxy : IBookService
{
    // lazy init: real service is only created when actually needed
    private RealBookService? _realService;
    private RealBookService Real => _realService ??= new RealBookService();

    // in-memory cache  key -> cached ServiceResponse
    private readonly Dictionary<string, ServiceResponse> _cache = new();
    private readonly List<ServiceLog> _proxyLogs = new();

    private void Log(string message, string type = "info")
    {
        _proxyLogs.Add(new ServiceLog
        {
            Timestamp = DateTime.Now.ToString("HH:mm:ss.fff"),
            Message   = message,
            Type      = type
        });
    }

    private ServiceResponse BuildCacheHit(string key)
    {
        var cached = _cache[key];
        var combined = new List<ServiceLog>(_proxyLogs);
        return new ServiceResponse
        {
            Books     = cached.Books,
            Logs      = combined,
            ElapsedMs = 0,   // sub-millisecond from cache
            FromCache = true
        };
    }

    private async Task<ServiceResponse> FetchAndCache(string key, Func<Task<ServiceResponse>> fetcher)
    {
        Log($"Cache MISS for key '{key}' — forwarding to RealBookService", "warn");
        Log("Delegating call to RealBookService...");

        var sw = Stopwatch.StartNew();
        var real = await fetcher();
        sw.Stop();

        _cache[key] = real;
        Log($"Response cached under key '{key}'", "cache");

        // merge proxy logs (pre-call) + real logs (the DB work)
        var combined = new List<ServiceLog>(_proxyLogs);
        combined.AddRange(real.Logs);

        return new ServiceResponse
        {
            Books     = real.Books,
            Logs      = combined,
            ElapsedMs = sw.ElapsedMilliseconds,
            FromCache = false
        };
    }

    // IBookService implementation
    public async Task<ServiceResponse> GetAllBooksAsync()
    {
        _proxyLogs.Clear();
        const string key = "all_books";

        Log("Proxy intercepted: GetAllBooksAsync()");
        Log($"Checking cache for key '{key}'...", "cache");

        if (_cache.ContainsKey(key))
        {
            Log($"Cache HIT for '{key}' — returning cached data instantly", "cache");
            return BuildCacheHit(key);
        }

        return await FetchAndCache(key, () => Real.GetAllBooksAsync());
    }

    public async Task<ServiceResponse> GetBookByIdAsync(int id)
    {
        _proxyLogs.Clear();
        string key = $"book_{id}";

        Log($"Proxy intercepted: GetBookByIdAsync({id})");
        Log($"Checking cache for key '{key}'...", "cache");

        if (_cache.ContainsKey(key))
        {
            Log($"Cache HIT for '{key}' — returning cached data instantly", "cache");
            return BuildCacheHit(key);
        }

        return await FetchAndCache(key, () => Real.GetBookByIdAsync(id));
    }

    public async Task<ServiceResponse> GetBooksByGenreAsync(string genre)
    {
        _proxyLogs.Clear();
        string key = $"genre_{genre.ToLower()}";

        Log($"Proxy intercepted: GetBooksByGenreAsync(\"{genre}\")");
        Log($"Checking cache for key '{key}'...", "cache");

        if (_cache.ContainsKey(key))
        {
            Log($"Cache HIT for '{key}' — returning cached data instantly", "cache");
            return BuildCacheHit(key);
        }

        return await FetchAndCache(key, () => Real.GetBooksByGenreAsync(genre));
    }

    public void ResetLogs()
    {
        _proxyLogs.Clear();
        _cache.Clear();   // also bust the cache so the demo is repeatable
    }
}
