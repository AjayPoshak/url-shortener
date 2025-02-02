db = db.getSiblingDB(process.env.MONGODB_DATABASE);

db.createUser({
  user: process.env.MONGODB_USER,
  pwd: process.env.MONGODB_PASSWORD,
  roles: [
    {
      role: "readWrite",
      db: process.env.MONGODB_DATABASE,
    },
  ],
});

// Create collections
db.createCollection("urls");
db.createCollection("users");
db.createCollection("analytics")

// Create indexes
db.urls.createIndex({ short_code: 1 }, { unique: true });
db.urls.createIndex({ created_at: 1 });

// Insert some initial data
db.urls.insertMany([
  {
    short_code: "example1",
    original_url: "https://example.com/very/long/url/1",
    created_at: new Date(),
    click_counter: 0,
  },
  {
    short_code: "example2",
    original_url: "https://example.com/very/long/url/2",
    created_at: new Date(),
    click_counter: 0,
  },
]);

// Create test database
db = db.getSiblingDB("myapp_test");
db.createCollection("urls");
db.createCollection("users");
