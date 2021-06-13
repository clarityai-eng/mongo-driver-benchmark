using System;
using System.Threading;
using System.Diagnostics;
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

            // Create a list of 100 documents
            var docList = new List<BsonDocument>();
            docList.Add(new BsonDocument{{ "title", "Dune" },{ "author", "Frank Herbert" }});

            for (int j = 0; j < 33; j++) {
                docList.Add(new BsonDocument{{ "title", "I, Robot" },{ "author", "Isaac Asimov" }});
                docList.Add(new BsonDocument{{ "title", "Foundation" },{ "author", "Isaac Asimov" }});
                docList.Add(new BsonDocument{{ "title", "Brave New World" },{ "author", "Aldous Huxley" }});
            }

            for (int i = 0; i < 6250; i++) {
              // Clone list to avoid duplicated ids 
              List<BsonDocument> cloned_list = docList.ConvertAll(book => new BsonDocument(book));
              collection.InsertMany(cloned_list);
            }
        }
    }

    class Program
    {
        static void Main(string[] args)
        {
            Stopwatch stopWatch = new Stopwatch();
            stopWatch.Start();

            var client = new MongoClient("mongodb://localhost:27017");
            var database = client.GetDatabase("test");
            var collection = database.GetCollection<BsonDocument>("dotnet");

            database.DropCollection("dotnet");

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

            stopWatch.Stop();
            Console.WriteLine("Insert time " + stopWatch.ElapsedMilliseconds);

            stopWatch.Restart();
            // Query documents
            var options = new FindOptions { BatchSize = 1000 };
            var filter = Builders<BsonDocument>.Filter.Eq("author", "Isaac Asimov");
            var cursor = collection.Find(filter, options).ToCursor();
            var totalChars = 0;
            foreach (var document in cursor.ToEnumerable())
            {
                var title = document.GetValue("title").ToString();
                var author = document.GetValue("author").ToString();
                totalChars += title.Length + author.Length;   
            }

            Console.WriteLine("Total Chars " + totalChars);

            stopWatch.Stop();
            Console.WriteLine("Find time " + stopWatch.ElapsedMilliseconds);
        }
    }
}
