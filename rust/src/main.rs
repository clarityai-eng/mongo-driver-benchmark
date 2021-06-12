use mongodb::{
    options::FindOptions,
    bson::{doc, Bson},
    sync::Client
};

use std::thread;
use std::time::{Instant};

const BATCH_SIZE:u32 = 1000;
const DB_NAME:&str = "test";
const COLLECTION_NAME:&str = "rust";

fn main() {
    let uri = "mongodb://localhost:27017";
    println!("Connecting to Mongo DB: {}", uri);

    let start = Instant::now();

    let client = Client::with_uri_str(uri).unwrap();
    let database = client.database(DB_NAME);
    let collection = database.collection(COLLECTION_NAME);

    collection.drop(None).unwrap();

    let mut threads = vec!();
    for _ in 0..8 {

        let mut docs = vec![
            doc! { "title": "Dune", "author": "Frank Herbert" },
        ];

        for _ in 0..33 {
            docs.push(doc! { "title": "I, Robot", "author": "Isaac Asimov" });
            docs.push(doc! { "title": "Foundation", "author": "Isaac Asimov" });
            docs.push(doc! { "title": "Brave New World", "author": "Aldous Huxley" });
        }

        threads.push(thread::spawn(move || {
            let client = Client::with_uri_str(uri).unwrap();
            let database = client.database(DB_NAME);
            let collection = database.collection(COLLECTION_NAME);
        
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

    println!("Insert time: {:?}", duration);

    // Now let's fetch all collection's data
    let start = Instant::now();

    // Query the documents in the collection with a filter and an option.
    let filter = doc! { "author": "Isaac Asimov" };
    let find_options = FindOptions::builder()
        .batch_size(BATCH_SIZE)
        .build();
    let cursor = collection.find(filter, find_options).unwrap();

    let mut total_chars = 0;
    // Iterate over the results of the cursor.
    for result in cursor {
        match result {
            Ok(document) => {
                let title = document.get("title").and_then(Bson::as_str).unwrap();
                let author = document.get("author").and_then(Bson::as_str).unwrap();
                total_chars = total_chars + title.len() + author.len();
            }
            Err(e) => println!("{:?}", e),
        } 
    }

    println!("Total chars: {:?}", total_chars);

    let duration = start.elapsed();
    println!("Fetch time: {:?}", duration);
}
