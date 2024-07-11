import base64
import openai
from dotenv import load_dotenv
import os
from flask import Flask, request, jsonify
from PIL import Image

# Load environment variables from .env file
load_dotenv()

app = Flask(__name__)

def convert_jpg_to_png(jpg_path, output_path):
    img = Image.open(jpg_path)
    img.save(output_path, 'PNG')
    return output_path

def describe_image(image_path, model="gpt-4o", temperature=0.0):
    # Get the API key from environment variables
    api_key = os.getenv("OPENAI_API_KEY")
    if api_key is None:
        raise ValueError("API key not found in environment variables")

    # Initialize OpenAI instance
    chatinstance = openai.OpenAI(api_key=api_key)

    # Function to encode image to base64
    def encode_image64(image_path):
        with open(image_path, "rb") as image_file:
            return base64.b64encode(image_file.read()).decode("utf-8")

    # Convert image to PNG if necessary
    if image_path.lower().endswith('.jpg') or image_path.lower().endswith('.jpeg'):
        image_path = convert_jpg_to_png(image_path, image_path.rsplit('.', 1)[0] + '.png')

    # Encode the image
    base64_image = encode_image64(image_path)

    # Make the API call to get the image description
    response = chatinstance.chat.completions.create(
        model=model,
        messages=[
            {"role": "system", "content": "You are a helpful assistant. Describe the image in one sentence."},
            {"role": "user", "content": [
                {"type": "text", "text": "Describe the image."},
                {"type": "image_url", "image_url": {"url": f"data:image/png;base64,{base64_image}"}
                }
            ]}
        ],
        temperature=temperature
    )

    # Return the image description
    return response.choices[0].message.content

@app.route('/describe-image', methods=['POST'])
def describe_image_endpoint():
    data = request.json
    image_path = data.get('image_path')

    if not image_path:
        return jsonify({"error": "Image path is required"}), 400

    try:
        description = describe_image(image_path)
        return jsonify({"description": description})
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True)
