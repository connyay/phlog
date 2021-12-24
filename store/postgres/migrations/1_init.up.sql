CREATE TABLE IF NOT EXISTS posts(
    id SERIAL NOT NULL,
    title TEXT NOT NULL,
    blobs TEXT[][]
);
