import string
from datasketch import MinHash
import asyncio
import moto_pb2
import moto_pb2_grpc
import grpc

K_SHINGLE = 3


def shingle(item: moto_pb2.ItemText):
    text = item.Text
    text = text.translate(str.maketrans('', '', string.punctuation))

    text_for_shingle = text.split(' ')

    shingle_set = []

    for i in range(len(text_for_shingle) - K_SHINGLE):
        shingle_set.append(''.join(text_for_shingle[i:i + K_SHINGLE]))

    return shingle_set


def compare(first_signature_input, second_signature_input):
    first_minhash = MinHash(hashvalues=first_signature_input)
    second_minhash = MinHash(hashvalues=second_signature_input)

    return first_minhash.jaccard(second_minhash)


class ItemServicer(moto_pb2_grpc.ItemServiceServicer):
    def GetSignature(self, request, context):
        shingle_set = shingle(request)
        mh = MinHash()
        for el in shingle_set:
            mh.update(el.encode('utf8'))
        s = mh.hashvalues
        return moto_pb2.ItemSignature(Signature=s)


async def serve() -> None:
    server = grpc.aio.server()
    moto_pb2_grpc.add_ItemServiceServicer_to_server(
        ItemServicer(), server)
    server.add_insecure_port('[::]:50051')
    await server.start()
    await server.wait_for_termination()


if __name__ == '__main__':
    asyncio.get_event_loop().run_until_complete(serve())
