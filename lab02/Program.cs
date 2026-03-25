using ProxyPatternLab.Services;

var builder = WebApplication.CreateBuilder(args);

// register both services
builder.Services.AddSingleton<RealBookService>();
builder.Services.AddSingleton<BookServiceProxy>();
builder.Services.AddCors(o => o.AddDefaultPolicy(p => p.AllowAnyOrigin().AllowAnyMethod().AllowAnyHeader()));

var app = builder.Build();
app.UseCors();
app.UseDefaultFiles();   // serves wwwroot/index.html at "/"
app.UseStaticFiles();


/// Without Pattern

var noProxy = app.MapGroup("/api/no-proxy");

noProxy.MapGet("/books", async (RealBookService svc) =>
    Results.Ok(await svc.GetAllBooksAsync()));

noProxy.MapGet("/books/{id:int}", async (int id, RealBookService svc) =>
    Results.Ok(await svc.GetBookByIdAsync(id)));

noProxy.MapGet("/books/genre/{genre}", async (string genre, RealBookService svc) =>
    Results.Ok(await svc.GetBooksByGenreAsync(genre)));

noProxy.MapPost("/reset", (RealBookService svc) => { svc.ResetLogs(); return Results.Ok(); });




/// With Pattern


var withProxy = app.MapGroup("/api/with-proxy");

withProxy.MapGet("/books", async (BookServiceProxy svc) =>
    Results.Ok(await svc.GetAllBooksAsync()));

withProxy.MapGet("/books/{id:int}", async (int id, BookServiceProxy svc) =>
    Results.Ok(await svc.GetBookByIdAsync(id)));

withProxy.MapGet("/books/genre/{genre}", async (string genre, BookServiceProxy svc) =>
    Results.Ok(await svc.GetBooksByGenreAsync(genre)));

withProxy.MapPost("/reset", (BookServiceProxy svc) => { svc.ResetLogs(); return Results.Ok(); });


// api info
app.MapGet("/api/info", () => Results.Ok(new
{
    pattern = "Proxy Design Pattern",
    participants = new[]
    {
        new { role = "Subject",     type = "IBookService",     description = "Common interface for RealSubject and Proxy" },
        new { role = "RealSubject", type = "RealBookService",  description = "Actual service; performs the expensive DB operations" },
        new { role = "Proxy",       type = "BookServiceProxy", description = "Intercepts calls; adds caching, logging, lazy init" },
        new { role = "Client",      type = "Web UI / API",     description = "Uses IBookService — unaware of which implementation it gets" }
    }
}));

app.Run();
