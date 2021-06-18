from elasticsearch import Elasticsearch
from shingle_minhash import compare

INDEX_NAME = 'items'
COUNT = 2
FIELD = 'signature'


def get_signature_from_es():
    es = Elasticsearch()

    res = es.search(index=INDEX_NAME, body={'size': COUNT, 'query': {'match_all': {}}})

    first_text_from_es = res['hits']['hits'][0]['_source'][FIELD]
    second_text_from_es = res['hits']['hits'][1]['_source'][FIELD]

    return first_text_from_es, second_text_from_es


first_signature, second_signature = get_signature_from_es()
j = compare(first_signature, second_signature)
print(j)
