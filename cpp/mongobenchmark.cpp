#include <cstdint>
#include <iostream>
#include <vector>
#include <thread>
#include <chrono>
// NOTE: You need to compile and install the mongocxx driver first. You can follow the official guide:
// http://mongocxx.org/mongocxx-v3/installation/ 
#include <bsoncxx/json.hpp>
#include <mongocxx/client.hpp>
#include <mongocxx/stdx.hpp>
#include <mongocxx/uri.hpp>
#include <mongocxx/instance.hpp>
#include <bsoncxx/builder/stream/helpers.hpp>
#include <bsoncxx/builder/stream/document.hpp>
#include <bsoncxx/builder/stream/array.hpp>

using bsoncxx::builder::stream::close_array;
using bsoncxx::builder::stream::close_document;
using bsoncxx::builder::stream::document;
using bsoncxx::builder::stream::finalize;
using bsoncxx::builder::stream::open_array;
using bsoncxx::builder::stream::open_document;

void worker_thread();

int main(int argc, char const *argv[])
{
	std::chrono::steady_clock::time_point begin = std::chrono::steady_clock::now();
	
	mongocxx::instance instance{};
	mongocxx::client client{mongocxx::uri{"mongodb://localhost:27017"}};

	// Get the collection
	mongocxx::database db = client["test"];
	mongocxx::collection collection = db["cpp"];

	// Clean collection
	collection.drop();

	std::vector<std::thread> threads;
	// Split work into 8 threads
	for (int i=0; i<8; ++i) {
    	threads.emplace_back(std::thread(worker_thread));
	}
	// Wait for all of them
	for (auto& th : threads) {
    	th.join();
	}

	std::chrono::steady_clock::time_point end = std::chrono::steady_clock::now();
	std::cout << "Insert time = " << std::chrono::duration_cast<std::chrono::milliseconds>(end - begin).count() << " ms" << std::endl;

	// Find from inserted docs
	mongocxx::options::find options;
    options.batch_size(1000);

	// Let's now fetch now all documents sequentially
	mongocxx::cursor cursor = collection.find(
		document{} << "author" << "Isaac Asimov" 
			<< bsoncxx::builder::stream::finalize, options);
	
	for(auto doc : cursor) {
		std::string title = doc["title"].get_utf8().value.to_string();
		std::string author = doc["author"].get_utf8().value.to_string();
	}

	return 0;
}

void worker_thread() {
	// Connect to DB
	mongocxx::client client{mongocxx::uri{"mongodb://localhost:27017"}};
	// Get the collection
	mongocxx::database db = client["test"];
	mongocxx::collection collection = db["cpp"];

	// Create the 4 books
	auto builder = bsoncxx::builder::stream::document{};
	bsoncxx::document::value dune = builder
		<< "title" << "Dune"
		<< "author" << "Frank Herbert"
		<< bsoncxx::builder::stream::finalize;

	bsoncxx::document::value i_robot = builder
		<< "title" << "I, Robot"
		<< "author" << "Isaac Asimov"
		<< bsoncxx::builder::stream::finalize;

	bsoncxx::document::value foundation = builder
		<< "title" << "Foundation"
		<< "author" << "Isaac Asimov"
		<< bsoncxx::builder::stream::finalize;

	bsoncxx::document::value brave_new_world = builder
		<< "title" << "Brave New World"
		<< "author" << "Aldous Huxley"
		<< bsoncxx::builder::stream::finalize;

	// Create array of books
	std::vector<bsoncxx::document::value> documents;
	documents.push_back(dune);
	for(int i = 0; i < 33; i++) {
		documents.push_back(i_robot);
		documents.push_back(foundation);
		documents.push_back(brave_new_world);
	}

	// Insert books in batches of 100
	for(int i=0; i < 6250; i++) {
		collection.insert_many(documents);
	}
}