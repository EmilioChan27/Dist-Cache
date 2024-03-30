SELECT Authors.*
FROM IWSchema.Articles
JOIN IWSchema.Authors ON Articles.AuthorId = Authors.Id
WHERE Articles.Id = 2;
GO