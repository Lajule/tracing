CREATE TABLE authors (
  id   BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name text    NOT NULL,
  bio  text
);

INSERT INTO authors
  (name, bio)
VALUES
  ('Douglas Adams', 'Douglas Noel Adams was an English author, humourist, and screenwriter, best known for The Hitchhiker''s Guide to the Galaxy. Originally a 1978 BBC radio comedy, The Hitchhiker''s Guide to the Galaxy developed into a "trilogy" of five books that sold more than 15 million copies in his lifetime.'),
  ('J. R. R. Tolkien', 'John Ronald Reuel Tolkien CBE FRSL was an English writer and philologist. He was the author of the high fantasy works The Hobbit and The Lord of the Rings. From 1925 to 1945, Tolkien was the Rawlinson and Bosworth Professor of Anglo-Saxon and a Fellow of Pembroke College, both at the University of Oxford.'),
  ('Frank Herbert', 'Franklin Patrick Herbert Jr. was an American science-fiction author, best known for his 1965 novel Dune and its five sequels. He also wrote short stories and worked as a newspaper journalist, photographer, book reviewer, ecological consultant, and lecturer.');
