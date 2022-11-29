CREATE TABLE
	labels (
		hash VARCHAR(64) NOT NULL,
		key VARCHAR NOT NULL,
		value VARCHAR NOT NULL,
		UNIQUE (hash, key, value),
		FOREIGN KEY (hash) REFERENCES blobs (hash) ON DELETE CASCADE
	);

CREATE INDEX
	labels_hash ON labels (hash);

CREATE INDEX
	labels_hash_key ON labels (hash, key);
