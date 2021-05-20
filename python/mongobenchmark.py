import pymongo, threading, time

def worker(count):
    myclient = pymongo.MongoClient("mongodb://localhost:27017/")
    mydb = myclient["test"]
    mycol = mydb["python"]

    for i in range(6250):

        docs = [{ "title": "Dune", "author": "Frank Herbert" }]
        for j in range(33):
            docs.append({ "title": "I, Robot", "author": "Isaac Asimov" })
            docs.append({ "title": "Foundation", "author": "Isaac Asimov" })
            docs.append({ "title": "Brave New World", "author": "Aldous Huxley" })

        mycol.insert_many(docs)
        
    return

myclient = pymongo.MongoClient("mongodb://localhost:27017/")
mydb = myclient["test"]
mycol = mydb["python"]

mycol.drop()

start = time.time()

# Insert all docs using 8 threads
threads = list()
for i in range(8):
    t = threading.Thread(target=worker, args=(i,))
    threads.append(t)
    t.start()

for t in threads:
    t.join()

end = time.time()
print(end - start)

# Now read all docs sequentially

start = time.time()
cursor = mycol.find({ "author": "Isaac Asimov" }, batch_size=1000)
for doc in cursor:
    title = doc["title"]
    author = doc["author"]   

end = time.time()
print(end - start)
