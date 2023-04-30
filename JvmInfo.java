public class JvmInfo {

    public static void main(String[] args) {
        System
            .getProperties()
            .keySet()
            .stream()
            .filter(String.class::isInstance)
            .map(String.class::cast)
            .filter(key -> key.startsWith("java."))
            .sorted()
            .forEach(JvmInfo::printProperty);
    }

    private static void printProperty(String property) {
        System.out.printf("%s=%s%n", property, System.getProperty(property));
    }
}
