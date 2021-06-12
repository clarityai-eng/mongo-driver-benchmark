const { MongoClient } = require("mongodb");
// Connection URI
const uri = "mongodb://localhost:27017";
// Create a new MongoClient
const client = new MongoClient(uri);

async function run() {
  try {
    var start = new Date();

    // Connect the client to the server
    await client.connect();
    // Clean collection first. If doesn't exist yet, catch the exception
    await client.db("test").collection("node").drop().catch(e => {})

    let collection = client.db("test").collection("node");

    // Insert all documents, doing batches of 8 concurrent async calls to the same client 
    for (i=0; i<6250; i++) {

      let insert_promises = [];
      
      for(t=0; t<8; t++) {
        // Doc list needs to be created each time so they are not considered the same docs by MongoClient
        // (otherwise they will be rejected by the driver due to duplication)
        let docs = [
          { "title": "Dune", "author": "Frank Herbert" },
        ];
      
        for (d=0; d<33; d++) {
          docs.push({ "title": "I, Robot", "author": "Isaac Asimov" });
          docs.push({ "title": "Foundation", "author": "Isaac Asimov" });
          docs.push({ "title": "Brave New World", "author": "Aldous Huxley" });
        }

        insert_promises.push(collection.insertMany(docs));
      }

      await Promise.all(insert_promises);
    }

    let end = new Date() - start
    console.info('Insert execution time: %dms', end)
    
    start = new Date();
    // Now let's find some of all inserted docs
    const findCursor = await collection.find({ author: "Isaac Asimov" }).batchSize(1000);

    let totalChars = 0; 
    await findCursor.forEach((doc) => {
      let author = doc.author;
      let title = doc.title;
      totalChars += author.length + title.length;
    });

    console.info('Total Chars: %d', totalChars)

    end = new Date() - start
    console.info('Fetch execution time: %dms', end)

  } finally {
    // Ensures that the client will close when you finish/error
    await client.close();
  }
}
run().catch(console.dir);