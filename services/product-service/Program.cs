using Ecommerce.ProductService.Data;
using Ecommerce.ProductService.Observability;
using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSingleton<RequestMetrics>();

var appPort = GetIntConfig(builder.Configuration, "APP_PORT", 8082);
builder.WebHost.ConfigureKestrel(options => options.ListenAnyIP(appPort));

var connectionString = BuildConnectionString(builder.Configuration);
builder.Services.AddDbContext<ProductDbContext>(options => options.UseNpgsql(connectionString));

var app = builder.Build();

await EnsureDatabaseAsync(app.Services);

app.Use(async (context, next) =>
{
    await next();

    var metrics = context.RequestServices.GetRequiredService<RequestMetrics>();
    metrics.RegisterRequest(context.Response.StatusCode);
});

app.MapGet("/health", async (ProductDbContext dbContext, CancellationToken cancellationToken) =>
{
    try
    {
        var canConnect = await dbContext.Database.CanConnectAsync(cancellationToken);
        return canConnect
            ? Results.Ok(new { status = "ok" })
            : Results.Json(new { status = "unhealthy" }, statusCode: StatusCodes.Status503ServiceUnavailable);
    }
    catch
    {
        return Results.Json(new { status = "unhealthy" }, statusCode: StatusCodes.Status503ServiceUnavailable);
    }
});

app.MapGet("/metrics", async (RequestMetrics metrics, ProductDbContext dbContext, CancellationToken cancellationToken) =>
{
    var productCount = await dbContext.Products.CountAsync(cancellationToken);
    return Results.Text(metrics.Render(productCount), "text/plain; version=0.0.4");
});

app.MapControllers();

await app.RunAsync();

static int GetIntConfig(IConfiguration configuration, string key, int fallback)
{
    var rawValue = configuration[key];
    return int.TryParse(rawValue, out var parsedValue) ? parsedValue : fallback;
}

static string BuildConnectionString(IConfiguration configuration)
{
    var host = configuration["DB_HOST"] ?? "localhost";
    var port = configuration["DB_PORT"] ?? "5432";
    var database = configuration["DB_NAME"] ?? "products_db";
    var user = configuration["DB_USER"] ?? "ecommerce";
    var password = configuration["DB_PASSWORD"] ?? "ecommerce123";
    var sslMode = configuration["DB_SSLMODE"] ?? "disable";

    return $"Host={host};Port={port};Database={database};Username={user};Password={password};SSL Mode={sslMode};";
}

static async Task EnsureDatabaseAsync(IServiceProvider serviceProvider)
{
    await using var scope = serviceProvider.CreateAsyncScope();
    var dbContext = scope.ServiceProvider.GetRequiredService<ProductDbContext>();
    await dbContext.Database.EnsureCreatedAsync();
}
