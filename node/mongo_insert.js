const { MongoClient } = require("mongodb");
// Connection URI
const uri =
  "mongodb://localhost:27017";
// Create a new MongoClient
const client = new MongoClient(uri);
async function run() {
  try {
    var start = new Date();

    // Connect the client to the server
    await client.connect();
    let collection = client.db("test").collection("perro");
    
    for (i=0; i<3750; i++) {
        const docs = [
            { "title": "1984", "author": "George Orwell" },
        ];
        
        for (j=0; j<33; j++) {
          docs.push({ "title": "1984", "author": "George Orwell" });
          docs.push({ "title": "Animal Farm", "author": "George Orwell" });
          docs.push({ "title": "The Great Gatsby", "author": "F. Scott Fitzgerald" });
        }

        await collection.insertMany(docs);
    
        if( i%10000 == 0 ) {
            console.log("Inserted " + i*100)
        }
    }

    var end = new Date() - start
    console.info('Execution time: %dms', end)
    
  } finally {
    // Ensures that the client will close when you finish/error
    await client.close();
  }
}
run().catch(console.dir);