public class Main {
    public static void main(String[] args) {
        Worker worker = new Worker("allen");
        System.out.println("name: " + worker.getName());

        worker.setName("bob");
        System.out.println("name: " + worker.getName());
    }
}
