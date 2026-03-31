using Ecommerce.ProductService.Models;
using Microsoft.EntityFrameworkCore;

namespace Ecommerce.ProductService.Data;

public sealed class ProductDbContext(DbContextOptions<ProductDbContext> options) : DbContext(options)
{
    public DbSet<Product> Products => Set<Product>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        modelBuilder.Entity<Product>(entity =>
        {
            entity.ToTable("products");
            entity.HasKey(x => x.Id);

            entity.Property(x => x.Name)
                .HasMaxLength(150)
                .IsRequired();

            entity.Property(x => x.Description)
                .HasMaxLength(1000)
                .IsRequired();

            entity.Property(x => x.Category)
                .HasMaxLength(80)
                .IsRequired();

            entity.Property(x => x.Price)
                .HasPrecision(18, 2)
                .IsRequired();

            entity.Property(x => x.Stock)
                .IsRequired();

            entity.Property(x => x.CreatedAtUtc)
                .IsRequired();

            entity.Property(x => x.UpdatedAtUtc)
                .IsRequired();

            entity.HasIndex(x => x.Category);
            entity.HasIndex(x => x.Name);
        });
    }
}
