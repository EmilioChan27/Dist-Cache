CREATE TABLE IWSchema.Articles (
    Id INT IDENTITY(1, 1) NOT NULL PRIMARY KEY,
    AuthorId INT FOREIGN KEY REFERENCES IWSchema.Authors(Id),
    Category NVARCHAR(50),
    Title NVARCHAR(150),
    ImageUrl NVARCHAR(MAX),
    Content NVARCHAR(MAX),
    Tags NVARCHAR(150),
    CreatedAt DATETIME,
    UpdatedAt DATETIME,
    Likes INT,
    Size INT
);
GO

INSERT INTO IWSchema.Articles (AuthorId, Category, Title, ImageUrl, Content, Tags, CreatedAt, UpdatedAt, Likes, Size)
VALUES (1, N'International', N'International Article 1 by 1', N'RandomImageUrl', N'Lorem Ipsum about some random global event that is happening currently somewhere around the world I want to make the content something but I dont want to waste too much time typing this all out and I want to make sure that I have enough time to get this in and I really hope that it will be ok in order to have the search function work. I will also say that this article has something to do with Asian business and economies and politics because I want to give it some tag at least', N'World,Asia,Business', GETDATE(), GETDATE(), 0, 25);
GO

SELECT * FROM IWSchema.Articles;
GO


