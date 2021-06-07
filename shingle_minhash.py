import string
from datasketch import MinHash
from elasticsearch import Elasticsearch

K_SHINGLE = 3
INDEX_NAME = 'items'
COUNT = 2
FIELD = 'text'


def shingle(text):
    text = text.translate(str.maketrans('', '', string.punctuation))

    text_for_shingle = text.split(' ')

    shingle_set = []

    for i in range(len(text_for_shingle) - K_SHINGLE):
        shingle_set.append(''.join(text_for_shingle[i:i + K_SHINGLE]))

    return shingle_set


def compare_with_minhash(first_set, second_set):
    first_minhash, second_minhash = MinHash(), MinHash()

    for i in first_set:
        first_minhash.update(i.encode('utf8'))

    for i in second_set:
        second_minhash.update(i.encode('utf8'))

    return first_minhash.jaccard(second_minhash)


def compare(first_text_input, second_text_input):
    first_set = shingle(first_text_input)
    second_set = shingle(second_text_input)

    return compare_with_minhash(first_set, second_set)


def get_texts_from_es():
    es = Elasticsearch()

    res = es.search(index=INDEX_NAME, body={'size': COUNT, 'query': {'match_all': {}}})

    first_text_from_es = res['hits']['hits'][0]['_source'][FIELD]
    second_text_from_es = res['hits']['hits'][1]['_source'][FIELD]

    return first_text_from_es, second_text_from_es


first_text, second_text = get_texts_from_es()
j = compare(first_text, second_text)
print(j)
