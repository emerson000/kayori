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

input_data = sys.stdin.read()
articles = json.loads(input_data)
nltk.download('stopwords', quiet=True)
nltk.download('punkt', quiet=True)
nltk.download('punkt_tab', quiet=True)

def preprocess_text(article):
    text = article['title']
    text = text.lower()  # Convert to lowercase
    text = text.strip()
    text = re.sub(r'\d+', '', text)  # Remove numbers
    text = re.sub(r'[^\w\s]', '', text)  # Remove punctuation
    text = re.sub(r'\s+', ' ', text)  # Remove extra spaces
    tokens = nltk.tokenize.word_tokenize(text)  # Tokenize
    tokens = [word for word in tokens if word.isalpha()]  # Keep only alphabetic tokens
    tokens = [word for word in tokens if word not in nltk.corpus.stopwords.words('english')]  # Remove stopwords
    return ' '.join(tokens)

processed_headlines = [preprocess_text(article) for article in articles]

ollama_client = ollama.Client(host='http://ollama:11434')

response = ollama_client.embed(model='nomic-embed-text', input=processed_headlines)

vectors = np.array(response['embeddings'])

# Reduce dimensionality of embeddings
pca = PCA(n_components=0.98)  # Adjust the number of components as needed
reduced_vectors = pca.fit_transform(vectors)

# Create a list of tuples (index, reduced_vector)
indexed_reduced_vectors = list(enumerate(reduced_vectors))

# Perform clustering on reduced vectors
dbscan = DBSCAN(eps=0.7, min_samples=2)  # Adjust eps for better clustering
clusters = dbscan.fit_predict(reduced_vectors)

# Display clusters
clustered_headlines = {}
for (i, _), cluster in zip(indexed_reduced_vectors, clusters):
    cluster_key = str(cluster)  # Convert cluster key to string
    if cluster_key not in clustered_headlines:
        clustered_headlines[cluster_key] = []
    clustered_headlines[cluster_key].append(articles[i])

print(json.dumps(clustered_headlines, indent=2))