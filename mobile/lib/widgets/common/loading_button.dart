import 'package:flutter/material.dart';
import '../../utils/constants.dart';

class LoadingButton extends StatelessWidget {
  final VoidCallback? onPressed;
  final bool isLoading;
  final String text;
  final Color? backgroundColor;
  final Color? foregroundColor;
  final double? height;
  final double? borderRadius;

  const LoadingButton({
    super.key,
    required this.onPressed,
    required this.isLoading,
    required this.text,
    this.backgroundColor,
    this.foregroundColor,
    this.height,
    this.borderRadius,
  });

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      height: height ?? 56,
      child: ElevatedButton(
        onPressed: isLoading ? null : onPressed,
        style: ElevatedButton.styleFrom(
          backgroundColor: backgroundColor ?? AppColors.primary,
          foregroundColor: foregroundColor ?? Colors.white,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(borderRadius ?? 28),
          ),
          elevation: isLoading ? 0 : 3,
          shadowColor: (backgroundColor ?? AppColors.primary).withValues(
            alpha: 0.3,
          ),
        ),
        child: isLoading
            ? const SizedBox(
                width: 24,
                height: 24,
                child: CircularProgressIndicator(
                  color: Colors.white,
                  strokeWidth: 2,
                ),
              )
            : Text(
                text,
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w600,
                ),
              ),
      ),
    );
  }
}
