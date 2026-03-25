using ProxyPatternLab.Models;

namespace ProxyPatternLab.Services;


public interface IBookService
{
    Task<ServiceResponse> GetAllBooksAsync();
    Task<ServiceResponse> GetBookByIdAsync(int id);
    Task<ServiceResponse> GetBooksByGenreAsync(string genre);
    void ResetLogs();
}
