from concurrent import futures
import grpc
import sys
import os

# Add the parent directory to the Python path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from python_proto.claude import claude_pb2
from python_proto.claude import claude_pb2_grpc
from ielts.ielts_service import IELTSService

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service = IELTSService()

    server.add_insecure_port('[::]:50052')
    print("Server starting on port 50052...")
    server.start()

    # Create a sample request
    sample_request = IELTSService.EvaluationRequestToClaude(
        student_response="This is a sample student response.",
        passage="This is a sample passage.",
        question="This is a sample question.",
        complex_sentences="These are sample complex sentences.",
        advanced_vocabulary="This is sample advanced vocabulary.",
        cohesive_devices="These are sample cohesive devices."
    )

    # Create a mock gRPC context
    class MockContext:
        def set_code(self, code):
            pass
        def set_details(self, details):
            pass

    # Execute EvaluateIELTS
    print("Executing EvaluateIELTS...")
    response = service.EvaluateIELTS(sample_request, MockContext())
    print("EvaluateIELTS Response:")
    print(response)

    server.wait_for_termination()

if __name__ == '__main__':
    serve()