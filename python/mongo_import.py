import pymongo, threading, time

def worker(count):
    myclient = pymongo.MongoClient("mongodb://localhost:27017/")

    mydb = myclient["test"]

    mycol = mydb["conejo"]

    for i in range(6250):

        docs = [{ "title": "1984", "author": "George Orwell" }]
        for j in range(33):
            docs.append({ "title": "1984", "author": "George Orwell" })
            docs.append({ "title": "Animal Farm", "author": "George Orwell" })
            docs.append({ "title": "The Great Gatsby", "author": "F. Scott Fitzgerald" })

        mycol.insert_many(docs)
        
    return

myclient = pymongo.MongoClient("mongodb://localhost:27017/")
mydb = myclient["test"]
mycol = mydb["conejo"]

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
documents_read = 0
cursor = mycol.find({}, batch_size=1000)
for doc in cursor:
    title = doc["title"]
    author = doc["author"]   
    documents_read = documents_read + 1

end = time.time()
print(end - start)
