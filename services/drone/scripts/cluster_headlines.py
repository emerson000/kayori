import os
os.environ["LOKY_MAX_CPU_COUNT"] = "4"
import nltk
import json
import sys
import re
import ollama
import numpy as np
from sklearn.cluster import DBSCAN
from sklearn.decomposition import PCA
from sklearn.preprocessing import normalize
from nltk.corpus import stopwords
import asyncio

OLLAMA_BATCH_SIZE = 50
MAX_CLUSTER_SIZE = 1000

input_data = sys.stdin.read()
articles = json.loads(input_data)
nltk.download('stopwords', quiet=True)
nltk.download('punkt', quiet=True)
nltk.download('punkt_tab', quiet=True)
stopwords.words('english')

def preprocess_text(article):
    text = article.get('title', '')
    text = text.lower().strip()
    text = re.sub(r'\d+', '', text)  # Remove numbers
    text = re.sub(r'[^\w\s]', '', text)  # Remove punctuation
    text = re.sub(r'\s+', ' ', text)  # Remove extra spaces
    tokens = nltk.tokenize.word_tokenize(text)
    tokens = [word for word in tokens if word.isalpha() and word not in {'the', 'a', 'an', 'is', 'are'}]
    return ' '.join(tokens)

async def process_headlines(articles):
    loop = asyncio.get_event_loop()
    tasks = [loop.run_in_executor(None, preprocess_text, article) for article in articles]
    return await asyncio.gather(*tasks)

async def embed_headlines(processed_headlines, ollama_client):
    all_embeddings = []
    for i in range(0, len(processed_headlines), OLLAMA_BATCH_SIZE):
        send_status_update(f"Processing batch {i // OLLAMA_BATCH_SIZE + 1} of {len(processed_headlines) // OLLAMA_BATCH_SIZE + 1}...")
        batch = processed_headlines[i:i + OLLAMA_BATCH_SIZE]
        response = ollama_client.embed(model='nomic-embed-text', input=batch)
        all_embeddings.extend(response['embeddings'])
    return all_embeddings

def perform_clustering(vectors):
    pca = PCA(n_components=0.95)  # Adjust the number of components as needed
    reduced_vectors = pca.fit_transform(vectors)
    dbscan = DBSCAN(eps=0.7, min_samples=2)  # Adjust eps for better clustering
    clusters = dbscan.fit_predict(reduced_vectors)
    return clusters, reduced_vectors

def build_cluster_output(clusters, reduced_vectors, articles, all_embeddings):
    indexed_reduced_vectors = list(enumerate(reduced_vectors))
    clustered_headlines = {}
    for (i, _), cluster in zip(indexed_reduced_vectors, clusters):
        cluster_key = str(cluster)  # Convert cluster key to string
        if cluster_key not in clustered_headlines:
            clustered_headlines[cluster_key] = {
                'articles': [],
                'centroid': None
            }
        article = articles[i]
        row = {
            'article': article,
            'embedding': all_embeddings[i]
        }
        clustered_headlines[cluster_key]['articles'].append(row)
        clustered_headlines[cluster_key]['centroid'] = np.mean([a['embedding'] for a in clustered_headlines[cluster_key]['articles']], axis=0).tolist()
    # Remove clusters that exceed the maximum size
    clustered_headlines = {k: v for k, v in clustered_headlines.items() if len(v) <= MAX_CLUSTER_SIZE}
    return clustered_headlines

async def main():
    send_status_update("Starting clustering process.")
    processed_headlines = await process_headlines(articles)
    send_status_update("Preprocessing completed. Starting embedding.")
    ollama_client = ollama.Client(host='http://ollama:11434')
    all_embeddings = await embed_headlines(processed_headlines, ollama_client)
    send_status_update("Embedding completed. Starting clustering.")
    vectors = normalize(np.array(all_embeddings))
    clusters, reduced_vectors = perform_clustering(vectors)
    send_status_update("Clustering completed. Building cluster output.")
    clustered_headlines = build_cluster_output(clusters, reduced_vectors, articles, all_embeddings)
    send_status_update("Clustering process completed.")
    print(json.dumps(clustered_headlines, indent=2))

def send_status_update(message):
    os.write(3, message.encode('utf-8'))

if __name__ == "__main__":
    asyncio.run(main())