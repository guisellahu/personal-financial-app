from flask import Flask, request, jsonify
from sqlalchemy import create_engine
import pandas as pd
from sklearn.linear_model import LinearRegression
import joblib
import os

app = Flask(__name__)

# Configuración de la conexión a la base de datos
DATABASE_URL = "postgresql://postgres:password@db:5432/mydb"
engine = create_engine(DATABASE_URL)

# Carga del modelo preentrenado si existe
MODEL_FILE = 'linear_regression_model.pkl'
if os.path.exists(MODEL_FILE):
    model = joblib.load(MODEL_FILE)
else:
    model = LinearRegression()

@app.route('/train', methods=['POST'])
def train_model():
    user_id = request.json.get('user_id')
    if not user_id:
        return jsonify({'error': 'User ID is required'}), 400

    # Utiliza el estilo de marcadores de posición de psycopg2
    query = "SELECT amount, date_part('day', created_at) as day FROM money_flows WHERE user_id = %s"
    df = pd.read_sql(query, engine, params=(user_id,))
    if df.empty:
        return jsonify({'error': 'No data found for training'}), 404

    X = df[['day']]
    y = df['amount']
    model.fit(X, y)
    joblib.dump(model, MODEL_FILE)
    return jsonify({'message': 'Model trained and saved successfully'})

@app.route('/predict', methods=['POST'])
def predict():
    if not os.path.exists(MODEL_FILE):
        return jsonify({'error': 'Model is not trained yet'}), 400

    # Suponemos que el JSON entrante tiene un mes y un año para predecir cada día del mes
    month = request.json.get('month')
    year = request.json.get('year')
    if not month or not year:
        return jsonify({'error': 'Month and year are required'}), 400

    # Generar rango de fechas para el mes especificado
    start_date = pd.Timestamp(year=year, month=month, day=1)
    end_date = start_date + pd.offsets.MonthEnd()

    # Crear DataFrame para días del mes
    days_in_month = pd.date_range(start=start_date, end=end_date)
    df = pd.DataFrame({'day': days_in_month.day, 'date': days_in_month})

    # Realizar predicciones
    predictions = model.predict(df[['day']])

    # Formatear predicciones
    df['prediction'] = predictions
    #df['prediction'] = df['prediction'].apply(lambda x: f"${int(x):,}")

    # Convertir DataFrame a JSON para la respuesta
    df['date'] = df['date'].dt.strftime('%Y-%m-%d')  # Formato de date YYY-MM-DD
    result = df[['date', 'prediction']].to_dict(orient='records')
    return jsonify({"predictions": result})

if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5000)
