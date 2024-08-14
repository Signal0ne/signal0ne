import logging
import torch
import json
from sentence_transformers import CrossEncoder

logger = logging.getLogger(__name__)
enc = CrossEncoder('cross-encoder/ms-marco-MiniLM-L-6-v2', default_activation_function=torch.nn.Sigmoid())

def log_occurrences(collectedLogs: list, comparedFields: list) -> list:
    unique_logs_list = []

    for log in collectedLogs:
        new = True
        incoming_log_object = {
            "count": 1
        }
        print("KEYS: ",dict(log).keys())
        for collectedLogKey in dict(log).keys():
            incoming_log_object[collectedLogKey] = log[collectedLogKey]

        for i, unique_log_object in enumerate(unique_logs_list):
            incoming_log_object_copy = {}
            unique_log_object_copy = {}
            for compareField in comparedFields:
                incoming_log_object_copy[compareField] = incoming_log_object[compareField]
                unique_log_object_copy[compareField] = unique_log_object[compareField]
            cos_sim = enc.predict([(json.dumps(incoming_log_object_copy), json.dumps(unique_log_object_copy))])
            if cos_sim[0] < 0.3:
                new = True
            else:
                new = False
                unique_logs_list[i]['count'] += 1
                break
        

        if new:
            unique_logs_list.append(incoming_log_object)

    return unique_logs_list

    