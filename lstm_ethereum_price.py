# LSTM for international airline passengers problem with regression framing
import numpy
import matplotlib.pyplot as plt
from pandas import read_csv
import math
from keras.models import Sequential
from keras.layers import Dense
from keras.layers import LSTM
from sklearn.preprocessing import MinMaxScaler
from sklearn.metrics import mean_squared_error

# convert an array of values into a dataset matrix
def create_dataset(dataset, look_back=1):
	dataX, dataY = [], []
	for i in range(len(dataset)-look_back):
		a = dataset[i:(i+look_back), 0]
		dataX.append(a)
		dataY.append(dataset[i + look_back, 0])
	return numpy.array(dataX), numpy.array(dataY)

def create_production_dataset(dataset, look_back=1):
	dataX = []
	for i in range(len(dataset) - look_back + 1):
		a = dataset[i:(i+look_back), 0]
		dataX.append(a)
	return numpy.array(dataX)

# fix random seed for reproducibility
numpy.random.seed(7)
# load the dataset
dataframe = read_csv('data/ethereum_price.csv', usecols=[1], engine='python')
dataset = dataframe.values
dataset = dataset.astype('double')
# normalize the dataset
scaler = MinMaxScaler(feature_range=(0, 1))
dataset = scaler.fit_transform(dataset)
# split into train and test sets
train_size = int(len(dataset) * 0.67)
test_size = len(dataset) - train_size
train, test = dataset[0:train_size,:], dataset[train_size:len(dataset),:]
# reshape into X=t and Y=t+1
look_back = 60
trainX, trainY = create_dataset(train, look_back)
testX, testY = create_dataset(test, look_back)
wholeX, wholeY = create_dataset(dataset, look_back)
# reshape input to be [samples, time steps, features]
trainX = numpy.reshape(trainX, (trainX.shape[0], 1, trainX.shape[1]))
testX = numpy.reshape(testX, (testX.shape[0], 1, testX.shape[1]))
wholeX = numpy.reshape(wholeX, (wholeX.shape[0], 1, wholeX.shape[1]))
# create and fit the LSTM network
model = Sequential()
model.add(LSTM(4, input_shape=(1, look_back)))
model.add(Dense(1))
model.compile(loss='mean_squared_error', optimizer='adam')
model.fit(wholeX, wholeY, epochs=100, batch_size=1, verbose=2)
# make predictions
trainPredict = model.predict(trainX)
testPredict = model.predict(testX)
wholePredict = model.predict(wholeX)
# invert predictions
trainPredict = scaler.inverse_transform(trainPredict)
trainY = scaler.inverse_transform([trainY])
testPredict = scaler.inverse_transform(testPredict)
testY = scaler.inverse_transform([testY])
wholePredict = scaler.inverse_transform(wholePredict)
wholeY = scaler.inverse_transform([wholeY])

# calculate root mean squared error
trainScore = math.sqrt(mean_squared_error(trainY[0], trainPredict[:,0]))
print('Train Score: %.2f RMSE' % (trainScore))
testScore = math.sqrt(mean_squared_error(testY[0], testPredict[:,0]))
print('Test Score: %.2f RMSE' % (testScore))
# shift train predictions for plotting
'''trainPredictPlot = numpy.empty_like(dataset)
trainPredictPlot[:, :] = numpy.nan
trainPredictPlot[look_back:len(trainPredict)+look_back, :] = trainPredict'''
# shift test predictions for plotting
'''testPredictPlot = numpy.empty_like(dataset)
testPredictPlot[:, :] = numpy.nan
testPredictPlot[len(trainPredict)+(look_back*2)+1:len(dataset)-1, :] = testPredict'''
# shift whole predictions for plotting
wholePredictPlot = numpy.empty_like(dataset)
wholePredictPlot[:, :] = numpy.nan
wholePredictPlot[look_back:len(dataset), :] = wholePredict

# print(testPredict)
# print('dataset')
# print(dataset)
# print(testX)

for i in range(0, 364):
	# print(dataset)
	input = create_production_dataset(dataset[-look_back:], look_back)
	# print(input)
	input = numpy.reshape(input, (input.shape[0], 1, input.shape[1]))
	prediction = model.predict(input)
	# print(prediction)
	dataset = numpy.concatenate((dataset,  prediction))

# plot baseline and predictions
plt.plot(scaler.inverse_transform(dataset))
# plt.plot(trainPredictPlot)
# plt.plot(testPredictPlot)
# plt.plot(wholePredictPlot, color="#FF0000")
plt.show()

coefs = scaler.inverse_transform(dataset)[-365:, 0]
f = open("data/coefficients_ethereum.dat" , "w")
print(dataset[-365:, 0])
for i in range (0, 365):
	print(str(coefs[i]))
	f.writelines(str(coefs[i]) + "\n")