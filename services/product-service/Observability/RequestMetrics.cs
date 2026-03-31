using System.Globalization;
using System.Text;

namespace Ecommerce.ProductService.Observability;

public sealed class RequestMetrics
{
    private long _httpRequestsTotal;
    private long _httpServerErrorsTotal;

    public void RegisterRequest(int statusCode)
    {
        Interlocked.Increment(ref _httpRequestsTotal);

        if (statusCode >= 500)
        {
            Interlocked.Increment(ref _httpServerErrorsTotal);
        }
    }

    public string Render(int productsTotal)
    {
        var requestsTotal = Interlocked.Read(ref _httpRequestsTotal);
        var serverErrorsTotal = Interlocked.Read(ref _httpServerErrorsTotal);

        var builder = new StringBuilder();
        builder.AppendLine("# HELP ecommerce_product_service_http_requests_total Total HTTP requests handled by product-service.");
        builder.AppendLine("# TYPE ecommerce_product_service_http_requests_total counter");
        builder.Append("ecommerce_product_service_http_requests_total ")
            .AppendLine(requestsTotal.ToString(CultureInfo.InvariantCulture));
        builder.AppendLine("# HELP ecommerce_product_service_http_server_errors_total Total HTTP responses with status >= 500.");
        builder.AppendLine("# TYPE ecommerce_product_service_http_server_errors_total counter");
        builder.Append("ecommerce_product_service_http_server_errors_total ")
            .AppendLine(serverErrorsTotal.ToString(CultureInfo.InvariantCulture));
        builder.AppendLine("# HELP ecommerce_product_service_products_total Current number of products in the catalog.");
        builder.AppendLine("# TYPE ecommerce_product_service_products_total gauge");
        builder.Append("ecommerce_product_service_products_total ")
            .AppendLine(productsTotal.ToString(CultureInfo.InvariantCulture));

        return builder.ToString();
    }
}
