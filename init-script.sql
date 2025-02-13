CREATE TABLE products (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(255) NOT NULL
);

INSERT INTO products (name)
SELECT 'Product ' + CAST(n AS NVARCHAR(255))
FROM (SELECT ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) AS n
      FROM master.dbo.spt_values) t
WHERE n <= 1000;
