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
import warnings

# Suppress specific RuntimeWarning related to libc not being found
warnings.filterwarnings("ignore", message="libc not found. The ctypes module in Python 3.12 is maybe too old for this OS.")

OLLAMA_BATCH_SIZE = 50
MAX_CLUSTER_SIZE = 1000

input_data = sys.stdin.read()
all_data = json.loads(input_data)
with open('/tmp/cluster_headlines.json', 'w') as f:
    json.dump(all_data, f, indent=2)
articles = all_data.get('articles', [])
existing_clusters = all_data.get('clusters', [])
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
    pca = PCA(n_components=0.98)  # Adjust the number of components as needed
    with warnings.catch_warnings():
        warnings.filterwarnings('error')
        try:
            reduced_vectors = pca.fit_transform(vectors)
        except RuntimeWarning as e:
            send_status_update(f"RuntimeWarning during PCA: {e}")
            reduced_vectors = vectors  # Fallback to original vectors if PCA fails

    dbscan = DBSCAN(eps=0.7, min_samples=2)  # Adjust eps for better clustering
    clusters = dbscan.fit_predict(reduced_vectors)
    return clusters, reduced_vectors

def perform_clustering_with_existing(vectors, existing_clusters, threshold=0.7):
    # Convert existing cluster centroids to numpy array
    existing_centroids = np.array([cluster['centroid'] for cluster in existing_clusters]).astype('float32')
    vectors = vectors.astype('float32')

    # Implement Cosine Similarity – Measures the cosine of the angle between two vectors.
    def cosine_similarity(vec1, vec2):
        dot_product = np.dot(vec1, vec2)
        norm_vec1 = np.linalg.norm(vec1)
        norm_vec2 = np.linalg.norm(vec2)
        return dot_product / (norm_vec1 * norm_vec2)

    clusters = []
    for vector in vectors:
        similarities = [cosine_similarity(vector, centroid) for centroid in existing_centroids]
        max_similarity = max(similarities)
        if max_similarity >= threshold:
            cluster_index = np.argmax(similarities)
            clusters.append(cluster_index)
        else:
            clusters.append(None)  # or any other value to indicate no cluster assignment

    return clusters, vectors

def build_cluster_output(clusters, reduced_vectors, articles, all_embeddings, existing_clusters=None):
    indexed_reduced_vectors = list(enumerate(reduced_vectors))
    clustered_headlines = {}
    for (i, _), cluster in zip(indexed_reduced_vectors, clusters):
        cluster_key = str(cluster)  # Convert cluster key to string
        if cluster_key not in clustered_headlines:
            if cluster is None:
                send_status_update(f"Cluster {cluster} is None, skipping.")
                continue
            clustered_headlines[cluster_key] = {
                'articles': [],
                'centroid': None,
                'id': None
            }
            if existing_clusters:
                send_status_update(f"Cluster {cluster} already exists, reusing its ID.")
                clustered_headlines[cluster_key]['id'] = existing_clusters[cluster]['id']
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


def send_status_update(message):
    try:
        os.write(3, message.encode('utf-8'))
    except OSError:
        sys.stderr.write(message + '\n')

async def cluster_headlines():
    send_status_update("Starting full clustering process.")
    processed_headlines = await process_headlines(articles)
    send_status_update("Preprocessing completed. Starting embedding.")
    ollama_client = ollama.Client(host='http://ollama:11434')
    all_embeddings = await embed_headlines(processed_headlines, ollama_client)
    send_status_update("Embedding completed. Starting clustering.")
    vectors = normalize(np.array(all_embeddings))
    if len(existing_clusters) == 0:
        send_status_update("Proceeding with clustering.")
        clusters, reduced_vectors = perform_clustering(vectors)
    else:
        send_status_update("Existing clusters provided. Using them for clustering.")
        clusters, reduced_vectors = perform_clustering_with_existing(vectors, existing_clusters)

    send_status_update("Clustering completed. Building cluster output.")
    clustered_headlines = build_cluster_output(clusters, reduced_vectors, articles, all_embeddings, existing_clusters if len(existing_clusters) > 0 else None)
    send_status_update("Clustering process completed.")
    print(json.dumps(clustered_headlines, indent=2))

async def main():
    await cluster_headlines()
    

def send_status_update(message):
    try:
        os.write(3, message.encode('utf-8'))
    except OSError:
        sys.stderr.write(message + '\n')

if __name__ == "__main__":
    asyncio.run(main())