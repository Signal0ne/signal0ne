import logging
import torch
from sentence_transformers import CrossEncoder

logger = logging.getLogger(__name__)
enc = CrossEncoder('cross-encoder/ms-marco-MiniLM-L-6-v2', default_activation_function=torch.nn.Sigmoid())

def contents_similarity(similarityCase: str, contents: list) -> list:
    contents_search_results = []

    for content in contents:
        cos_sim = enc.predict([(similarityCase, content)])
        if cos_sim[0] > 0.3:
            contents_search_results.append(content)

    return contents_search_results
