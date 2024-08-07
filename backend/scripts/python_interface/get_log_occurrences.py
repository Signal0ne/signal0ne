import logging
import torch
import json
from sentence_transformers import CrossEncoder

logger = logging.getLogger(__name__)
enc = CrossEncoder('cross-encoder/ms-marco-MiniLM-L-6-v2', default_activation_function=torch.nn.Sigmoid())

def log_occurrences(collectedLogs: list):
    # What are the fields that we should do the comparison by??? should user decide on that???
    unique_logs_list = []

    for log in collectedLogs:
        new = True
        incoming_log_object = {
            "count": 1
        }
        for collectedLogKey in dict(log).keys():
            incoming_log_object[collectedLogKey] = log[collectedLogKey]

        for i, unique_log_object in enumerate(unique_logs_list):
            cos_sim = enc.predict([(json.dumps(incoming_log_object), json.dumps(unique_log_object))])
            if cos_sim[0] < 0.3:
                new = True
            else:
                new = False
                unique_logs_list[i]['count'] += 1
                break
        

        if new:
            unique_logs_list.append(incoming_log_object)
    