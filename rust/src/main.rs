use mongodb::{
    options::FindOptions,
    bson::{doc, Bson},
    sync::Client
};

use std::thread;
use std::time::{Instant};
//use std::env;

const BATCH_SIZE:u32 = 1000;

fn main() {
    let uri = "mongodb://localhost:27017";
    println!("Connecting to Mongo DB: {}", uri);

    // Clean BD
    Client::with_uri_str(uri).unwrap().database("test").collection("pollo").drop(None).unwrap();

    let start = Instant::now();

    let client = Client::with_uri_str(uri).unwrap();
    let database = client.database("test");
    let collection = database.collection("pollo");

    collection.drop(None).unwrap();

    let mut threads = vec!();
    for _ in 0..8 {

        threads.push(thread::spawn(move || {
            let client = Client::with_uri_str(uri).unwrap();
            let database = client.database("test");
            let collection = database.collection("pollo");

            let mut docs = vec![
                doc! { "title": "1984", "author": "George Orwell" },
            ];

            for _ in 0..33 {
                docs.push(doc! { "title": "1984", "author": "George Orwell" });
                docs.push(doc! { "title": "Animal Farm", "author": "George Orwell" });
                docs.push(doc! { "title": "The Great Gatsby", "author": "F. Scott Fitzgerald" });
            }
        
            for _ in 0..6250 {
                collection.insert_many(docs.clone(), None).unwrap();
            }

        }));
    }

    // Wait until all threads finish
    for thread in threads {
        let _ = thread.join();
    }

    let duration = start.elapsed();

    println!("Time spent in insert: {:?}", duration);

    // Now let's fetch all collection's data

    let client = Client::with_uri_str(uri).unwrap();
    let database = client.database("test");
    let collection = database.collection("pollo");

    let start = Instant::now();

    //let args: Vec<String> = env::args().collect();
    //let batch_size:u32 = args[1].parse().unwrap();

    // Query the documents in the collection with a filter and an option.
    let filter = doc! { "author": "George Orwell" };
    let find_options = FindOptions::builder()
        .batch_size(BATCH_SIZE)
        .build();
    let cursor = collection.find(filter, find_options).unwrap();

    let mut processed_docs = 0;

    // Iterate over the results of the cursor.
    for result in cursor {
        match result {
            Ok(document) => {
                let _title = document.get("title").and_then(Bson::as_str);
                let _author = document.get("author").and_then(Bson::as_str);
                processed_docs = processed_docs + 1;
            }
            Err(e) => println!("{:?}", e),
        } 
    }

    let duration = start.elapsed();
    println!("Time spent in processing {} documents: {:?}", processed_docs, duration);
}
