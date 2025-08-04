class ApiResponse<T> {
  final T? data;
  final String message;
  final bool isSuccess;

  ApiResponse.success(this.data) : message = '', isSuccess = true;

  ApiResponse.error(this.message) : data = null, isSuccess = false;
}
