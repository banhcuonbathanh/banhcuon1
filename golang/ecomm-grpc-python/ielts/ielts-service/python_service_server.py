import grpc
from concurrent import futures
import python_service_pb2
import python_service_pb2_grpc

class PythonServiceServicer(python_service_pb2_grpc.PythonServiceServicer):
    def ProcessData(self, request, context):
        input_data = request.input_data
        # Process the data (replace this with your actual logic)
        processed_data = f"Processed: {input_data}"
        return python_service_pb2.DataResponse(processed_data=processed_data)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    python_service_pb2_grpc.add_PythonServiceServicer_to_server(PythonServiceServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()