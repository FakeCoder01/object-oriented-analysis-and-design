import javax.swing.*;
import javax.swing.table.DefaultTableModel;
import java.awt.*;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.util.HashMap;
import java.util.Map;

class Record {
    private String key;
    private String value;

    public Record(String key, String value) {
        this.key = key;
        this.value = value;
    }
    public String getKey() { return key; }
    public String getValue() { return value; }
    public void setValue(String value) { this.value = value; }
}

class Gateway {
    private String host;
    private int port;
    private Socket socket;
    private OutputStream out;
    private InputStream in;

    public Gateway(String host, int port) {
        this.host = host;
        this.port = port;
    }

    public void connect() throws Exception {
        socket = new Socket(host, port);
        out = socket.getOutputStream();
        in = socket.getInputStream();
    }

    public void disconnect() throws Exception {
        if (socket != null) socket.close();
    }

    public String executeCommand(String cmd) throws Exception {
        out.write(cmd.getBytes());
        out.flush();

        byte[] buffer = new byte[1024];
        int bytesRead = in.read(buffer);
        if (bytesRead > 0) {
            return new String(buffer, 0, bytesRead).trim();
        }
        return null;
    }
}

class IdentityMap {
    private Map<String, Record> cache = new HashMap<>();
    private Gateway gateway;
    private JTextArea logger;

    public IdentityMap(Gateway gateway, JTextArea logger) {
        this.gateway = gateway;
        this.logger = logger;
    }

    public Record find(String key) {
        if (cache.containsKey(key)) {
            logger.append("[IDENTITY MAP] CACHE HIT! Returning existing object for key: " + key + "\n");
            return cache.get(key);
        }

        logger.append("[IDENTITY MAP] CACHE MISS. Fetching from MeowDB server...\n");
        try {
            String response = gateway.executeCommand("select " + key);
            if (response != null && !response.contains("not found")) {
                Record record = new Record(key, response);
                cache.put(key, record);
                logger.append("[IDENTITY MAP] Object instantiated and mapped.\n");
                return record;
            }
        } catch (Exception e) {
            logger.append("Error fetching: " + e.getMessage() + "\n");
        }
        return null;
    }

    public void save(String command, String key, String value) {
        try {
            String response = gateway.executeCommand(command + " " + key + " " + value);
            if (response != null && !response.contains("not found")) {
                if (cache.containsKey(key)) {
                    cache.get(key).setValue(value);
                    logger.append("[IDENTITY MAP] Existing object updated in map.\n");
                } else {
                    cache.put(key, new Record(key, value));
                    logger.append("[IDENTITY MAP] New object added to map.\n");
                }
            }
        } catch (Exception e) {
            logger.append("Error saving: " + e.getMessage() + "\n");
        }
    }

    public void delete(String key) {
        try {
            String response = gateway.executeCommand("delete " + key);
            if (response != null && !response.contains("not found")) {
                cache.remove(key);
                logger.append("[IDENTITY MAP] Object removed from map.\n");
            } else {
                logger.append("[DATABASE] Key not found to delete.\n");
            }
        } catch (Exception e) {
            logger.append("Error deleting: " + e.getMessage() + "\n");
        }
    }

    public Map<String, Record> getAllCached() {
        return cache;
    }
}

public class Client extends JFrame {
    private IdentityMap identityMap;
    private Gateway gateway;
    private JTextArea consoleArea;
    private DefaultTableModel tableModel;
    private JTextField keyField, valueField, hostField;
    private JButton connectBtn;
    private boolean connected = false;

    public Client() {
        setTitle("MeowDB - Identity Map Pattern");
        setSize(850, 500);
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        setLayout(new BorderLayout());

        JPanel top = new JPanel();
        hostField = new JTextField("meowdb", 10);
        connectBtn = new JButton("Connect via TCP");
        top.add(new JLabel("Server Host:")); top.add(hostField);
        top.add(new JLabel("Port: 6969")); top.add(connectBtn);
        add(top, BorderLayout.NORTH);

        tableModel = new DefaultTableModel(new String[]{"Memory Address (Fake)", "Key", "Value"}, 0);
        JTable table = new JTable(tableModel);
        add(new JScrollPane(table), BorderLayout.CENTER);

        JPanel left = new JPanel(new GridLayout(6, 1, 5, 5));
        left.setBorder(BorderFactory.createTitledBorder("Commands"));
        keyField = new JTextField();
        valueField = new JTextField();
        left.add(new JLabel("Key:"));
        left.add(keyField);
        left.add(new JLabel("Value:"));
        left.add(valueField);

        JButton btnSelect = new JButton("Select");
        JButton btnInsert = new JButton("Insert / Update");
        JButton btnDelete = new JButton("Delete");
        left.add(btnSelect);
        left.add(btnInsert);
        left.add(btnDelete);
        add(left, BorderLayout.WEST);

        consoleArea = new JTextArea(8, 50);
        consoleArea.setBackground(Color.BLACK);
        consoleArea.setForeground(Color.GREEN);
        consoleArea.setFont(new Font("Monospaced", Font.PLAIN, 12));
        add(new JScrollPane(consoleArea), BorderLayout.SOUTH);

        connectBtn.addActionListener(e -> toggleConnection());
        btnSelect.addActionListener(e -> {
            Record r = identityMap.find(keyField.getText());
            if (r != null) {
                valueField.setText(r.getValue());
                refreshTable();
            } else {
                consoleArea.append("[DATABASE] Key not found.\n");
            }
        });

        btnInsert.addActionListener(e -> {
            identityMap.save("insert", keyField.getText(), valueField.getText());
            refreshTable();
        });

        btnDelete.addActionListener(e -> {
            identityMap.delete(keyField.getText());
            keyField.setText(""); valueField.setText("");
            refreshTable();
        });
    }

    private void toggleConnection() {
        if (!connected) {
            try {
                gateway = new Gateway(hostField.getText(), 6969);
                gateway.connect();
                identityMap = new IdentityMap(gateway, consoleArea);
                consoleArea.append("Connected to MeowDB.\n");
                connected = true;
                connectBtn.setText("Disconnect");
            } catch (Exception ex) {
                consoleArea.append("Connection failed: " + ex.getMessage() + "\n");
            }
        } else {
            try { gateway.disconnect(); } catch (Exception ex) {}
            connected = false;
            connectBtn.setText("Connect via TCP");
            consoleArea.append("Disconnected.\n");
        }
    }

    private void refreshTable() {
        tableModel.setRowCount(0);
        for (Record rec : identityMap.getAllCached().values()) {
            tableModel.addRow(new Object[]{"@" + Integer.toHexString(System.identityHashCode(rec)), rec.getKey(), rec.getValue()});
        }
    }

    public static void main(String[] args) {
        SwingUtilities.invokeLater(() -> new Client().setVisible(true));
    }
}
