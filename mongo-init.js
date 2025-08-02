// MongoDB initialization script
// This script will be executed when MongoDB container starts

// Switch to the taskmanager database
db = db.getSiblingDB('taskmanager');

// Create collections with validation (optional)
db.createCollection('users', {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: [ "username", "email", "password" ],
         properties: {
            username: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            email: {
               bsonType: "string",
               pattern: "^.+@.+\..+$",
               description: "must be a valid email address and is required"
            },
            password: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            first_name: {
               bsonType: "string",
               description: "must be a string"
            },
            last_name: {
               bsonType: "string",
               description: "must be a string"
            }
         }
      }
   }
});

db.createCollection('tasks', {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: [ "title", "user_id" ],
         properties: {
            title: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            description: {
               bsonType: "string",
               description: "must be a string"
            },
            status: {
               enum: [ "pending", "in_progress", "completed", "cancelled" ],
               description: "can only be one of the enum values"
            },
            priority: {
               enum: [ "low", "medium", "high", "urgent" ],
               description: "can only be one of the enum values"
            },
            user_id: {
               bsonType: "string",
               description: "must be a string and is required"
            }
         }
      }
   }
});

db.createCollection('categories');

// Create indexes for better performance
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "username": 1 }, { unique: true });
db.tasks.createIndex({ "user_id": 1 });
db.tasks.createIndex({ "status": 1 });
db.tasks.createIndex({ "priority": 1 });
db.tasks.createIndex({ "due_date": 1 });
db.categories.createIndex({ "user_id": 1 });

// Create a sample admin user (optional)
// Note: In production, you should hash the password properly
/*
db.users.insertOne({
    _id: ObjectId(),
    username: "admin",
    email: "admin@taskmanager.com",
    password: "$2a$10$example_hashed_password", // This should be properly hashed
    first_name: "Admin",
    last_name: "User",
    created_at: new Date(),
    updated_at: new Date()
});
*/

print("MongoDB initialization completed for taskmanager database");
