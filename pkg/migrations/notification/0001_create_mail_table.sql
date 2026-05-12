CREATE TABLE IF NOT EXISTS mail (
                                    id SERIAL PRIMARY KEY,
                                    sender TEXT NOT NULL,
                                    receiver TEXT NOT NULL,
                                    message TEXT NOT NULL,
                                    created_at TIMESTAMP NOT NULL
);