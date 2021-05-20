using System.Threading;
using System.Collections.Generic;
using MongoDB.Bson;
using MongoDB.Driver;

namespace mongobenchmark
{
    public class WorkerThread {
        public static void insertDocs()
        {
            var client = new MongoClient("mongodb://localhost:27017");

            var database = client.GetDatabase("test");
            var collection = database.GetCollection<BsonDocument>("dotnet");

            for (int i = 0; i < 6250; i++) {
                // Create a list of 100 documents
                var docList = new List<BsonDocument>();
                docList.Add(new BsonDocument{{ "title", "Dune" },{ "author", "Frank Herbert" }});

                for (int j = 0; j < 33; j++) {
                    docList.Add(new BsonDocument{{ "title", "I, Robot" },{ "author", "Isaac Asimov" }});
                    docList.Add(new BsonDocument{{ "title", "Foundation" },{ "author", "Isaac Asimov" }});
                    docList.Add(new BsonDocument{{ "title", "Brave New World" },{ "author", "Aldous Huxley" }});
                }

              collection.InsertMany(docList);
            }
        }
    }

    class Program
    {
        static void Main(string[] args)
        {
            var client = new MongoClient("mongodb://localhost:27017");
            var database = client.GetDatabase("test");
            var collection = database.GetCollection<BsonDocument>("dotnet");

            // Create and start all worker threads 
            var threads = new List<Thread>();
            for(int i=0; i<8; i++){
                var newThread = new Thread(WorkerThread.insertDocs);
                threads.Add(newThread);
                newThread.Start();
            }
            
            // Wait for all threads to finish
            foreach(Thread thread in threads) {
                thread.Join();
            }

            // Query documents
            var options = new FindOptions<BsonDocument> { BatchSize = 1000 };
            var filter = Builders<BsonDocument>.Filter.Eq("author", "Isaac Asimov");
            var cursor = collection.Find(filter).ToCursor();
            foreach (var document in cursor.ToEnumerable())
            {
                var title = document.GetValue("title").ToString();
                var author = document.GetValue("author").ToString();   
            }
        }
    }
}
