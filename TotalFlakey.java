import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.DocumentBuilder;
import org.w3c.dom.Document;
import org.w3c.dom.NodeList;
import org.w3c.dom.Node;
import org.w3c.dom.NamedNodeMap;
import org.w3c.dom.Element;
import java.io.File;
import java.util.HashMap;
import java.util.Map;
import java.io.FileWriter;
import java.io.BufferedWriter;
import java.io.IOException;
import javax.xml.parsers.ParserConfigurationException;
import org.xml.sax.SAXException;
import java.io.FilenameFilter; 
import javax.xml.transform.Transformer;
import javax.xml.transform.TransformerException;
import javax.xml.transform.TransformerFactory;
import javax.xml.transform.dom.DOMSource;
import javax.xml.transform.stream.StreamResult;
import org.w3c.dom.Attr;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.Calendar;
import java.util.Date;
import javax.xml.transform.OutputKeys;
import java.nio.file.Paths;
import java.net.URI;
import org.xml.sax.InputSource;
import java.io.StringReader;
import com.google.cloud.ServiceOptions;
import com.google.cloud.storage.Blob;
import com.google.cloud.storage.BlobId;
import com.google.cloud.storage.Storage.BlobListOption;
import com.google.cloud.storage.Storage;
import com.google.cloud.storage.StorageOptions;
import com.google.api.gax.paging.Page;

// project name: istioFlakeyTest in gcloud
// compile: javac -cp ".:google-cloud-storage-1.74.0.jar" TotalFlakey.java
// run: java -cp ".:jars/*" TotalFlakey

public class TotalFlakey {

	private static HashMap<String, Pair<Integer, Integer>> addSuccessfulCase(HashMap<String, Pair<Integer, Integer>> caseCollection, String caseName) {
		if (caseCollection.containsKey(caseName)) {
	    	Pair<Integer, Integer> caseResult = caseCollection.get(caseName);
	    	caseResult.setSecond(caseResult.getSecond() + 1);
	    	caseCollection.put(caseName, caseResult);
	    } else {
	    	Pair<Integer, Integer> caseResult = new Pair<Integer, Integer> (0, 1);
	    	caseCollection.put(caseName, caseResult);
	    }
	    return caseCollection;
	}

	public static void identifyFailures(HashMap<String, Pair<Pair<Integer, Integer>, HashMap<String, Pair<Integer, Integer>>>> flakey, Document doc) {
		int tests;
		NodeList nodeList = doc.getElementsByTagName("testsuite");
	    for(int x=0,size= nodeList.getLength(); x<size; x++) {
	    	Node curNode = nodeList.item(x);
	    	
	    	if (curNode.getNodeType() == Node.ELEMENT_NODE) {

	    		if (curNode.hasAttributes()) {
	    			NamedNodeMap nodeMap = curNode.getAttributes();
	    			String suiteName = nodeMap.getNamedItem("name").getNodeValue();
	    			int numSuiteFailures = Integer.parseInt(nodeMap.getNamedItem("failures").getNodeValue());
	    			int numSuiteTests = Integer.parseInt(nodeMap.getNamedItem("tests").getNodeValue());

	    			if (flakey.containsKey(suiteName)) {
	    				Pair<Pair<Integer, Integer>, HashMap<String, Pair<Integer, Integer>>> result = flakey.get(suiteName);
	    				Pair<Integer, Integer> suiteResult = result.getFirst();
	    				HashMap<String, Pair<Integer, Integer>> caseCollection = result.getSecond();
	    				int suiteTotal = suiteResult.getSecond();
	    				suiteResult.setSecond(suiteTotal + 1);

	    				if (numSuiteFailures != 0) {
							int suiteFailure = suiteResult.getFirst();
							suiteResult.setFirst(suiteFailure + 1);
							if (curNode.hasChildNodes()) {
	    						NodeList childNodes = curNode.getChildNodes();
	    						for (int y = 0; y < childNodes.getLength(); y ++) {
	    							Node testCase = childNodes.item(y);
	    							if (testCase.getNodeType() == Node.ELEMENT_NODE && testCase.getNodeName().equals("testcase")) {
	    								NamedNodeMap casemap = testCase.getAttributes();
	    								String caseName = suiteName + "/" + casemap.getNamedItem("name").getNodeValue();
	    								
	    								NodeList caseChildren = testCase.getChildNodes();
    									Boolean containsFailure = false;
    									for (int k = 0; k < caseChildren.getLength(); k ++) {
    										Node child = caseChildren.item(k);
    										if (child.getNodeName().equals("failure")) {
    											containsFailure = true;
    											if (caseCollection.containsKey(caseName)) {
    												Pair<Integer, Integer> caseResult = caseCollection.get(caseName);
    												caseResult.setFirst(caseResult.getFirst() + 1);
    												caseResult.setSecond(caseResult.getSecond() + 1);
    											
    												
    												caseCollection.put(caseName, caseResult);
    											} else {
    												Pair<Integer, Integer> caseResult = new Pair<Integer, Integer>(1, 1);
    												
    												caseCollection.put(caseName, caseResult);

    											}
    											break;
    										}
    									}
    									if (containsFailure == false) {
    										caseCollection = addSuccessfulCase(caseCollection, caseName);
    									}
    								}
    							}
    						}
						} else {
							if (curNode.hasChildNodes()) {
	    						NodeList childNodes = curNode.getChildNodes();
	    						for (int y = 0; y < childNodes.getLength(); y ++) {
	    							Node testCase = childNodes.item(y);
	    							if (testCase.getNodeType() == Node.ELEMENT_NODE && testCase.getNodeName().equals("testcase")) {
	    								NamedNodeMap casemap = testCase.getAttributes();
	    								String caseName = suiteName + "/" + casemap.getNamedItem("name").getNodeValue();
										caseCollection = addSuccessfulCase(caseCollection, caseName);
									}
								}
							}
						}
						result.setFirst(suiteResult);
						result.setSecond(caseCollection);
						flakey.put(suiteName, result);

					} else {
						Pair<Integer, Integer> suiteResult = new Pair<>(0, 1);
						HashMap<String, Pair<Integer, Integer>> caseCollection = new HashMap<>();
						if (numSuiteFailures != 0) {
							int suiteFailure = suiteResult.getFirst();
							suiteResult.setFirst(suiteFailure + 1);
							if (curNode.hasChildNodes()) {
	    						NodeList childNodes = curNode.getChildNodes();
	    						for (int y = 0; y < childNodes.getLength(); y ++) {
	    							Node testCase = childNodes.item(y);
	    							if (testCase.getNodeType() == Node.ELEMENT_NODE && testCase.getNodeName().equals("testcase")) {
	    								NamedNodeMap casemap = testCase.getAttributes();
	    								String caseName = suiteName + "/" + casemap.getNamedItem("name").getNodeValue();
	    								
	    								NodeList caseChildren = testCase.getChildNodes();
    									Boolean containsFailure = false;
    									for (int k = 0; k < caseChildren.getLength(); k ++) {
    										Node child = caseChildren.item(k);
    										if (child.getNodeName().equals("failure")) {
    											containsFailure = true;
    											Pair<Integer, Integer> caseResult = new Pair<Integer, Integer>(1, 1);
    												
    											caseCollection.put(caseName, caseResult);
    										}
    										break;
    									}
    									if (containsFailure == false) {
    										caseCollection = addSuccessfulCase(caseCollection, caseName);
    									}
    								}
    							}
    						}
    					} else {
							if (curNode.hasChildNodes()) {
	    						NodeList childNodes = curNode.getChildNodes();
	    						for (int y = 0; y < childNodes.getLength(); y ++) {
	    							Node testCase = childNodes.item(y);
	    							if (testCase.getNodeType() == Node.ELEMENT_NODE && testCase.getNodeName().equals("testcase")) {
	    								NamedNodeMap casemap = testCase.getAttributes();
	    								String caseName = suiteName + "/" + casemap.getNamedItem("name").getNodeValue();
										caseCollection = addSuccessfulCase(caseCollection, caseName);
									}
								}
							}
						}
						
						Pair<Pair<Integer, Integer>, HashMap<String, Pair<Integer, Integer>>> result = new Pair<>(suiteResult, caseCollection);
						flakey.put(suiteName, result);
					}
				}
			}
		}
	}

	private static void printFlakey(HashMap<String, Pair<Pair<Integer, Integer>, HashMap<String, Pair<Integer, Integer>>>> flakey, String filePath) throws TransformerException, ParserConfigurationException{

		String xmlPattern = "/^[a-zA-Z_:][a-zA-Z0-9\\.\\-_:]*$/";
		Pattern pattern = Pattern.compile(xmlPattern);


		DocumentBuilderFactory documentFactory = DocumentBuilderFactory.newInstance();
 
        DocumentBuilder documentBuilder = documentFactory.newDocumentBuilder();

        Document document = documentBuilder.newDocument();

        Element root = document.createElement("testsuites");
        document.appendChild(root);

        for (String suiteName : flakey.keySet()) {

        	Pair<Pair<Integer, Integer>, HashMap<String, Pair<Integer, Integer>>> result = flakey.get(suiteName);
        	Pair<Integer, Integer> suiteResult = result.getFirst();
        	HashMap<String, Pair<Integer, Integer>> caseCollection = result.getSecond();
        	Element testsuite = document.createElement("testsuite");
        	Attr attrName = document.createAttribute("name");
        	attrName.setValue(suiteName);
        	testsuite.setAttributeNode(attrName);
        	//Element testsuite = document.createElement(suiteName);

        	Attr suiteFailure = document.createAttribute("failures");
            suiteFailure.setValue(Integer.toString(suiteResult.getFirst()));
            testsuite.setAttributeNode(suiteFailure);

            Attr suiteTotal = document.createAttribute("total");
            suiteTotal.setValue(Integer.toString(suiteResult.getSecond()));
            testsuite.setAttributeNode(suiteTotal);


            for (String caseName : caseCollection.keySet()) {
            	Pair<Integer, Integer> caseResult = caseCollection.get(caseName);
            	Element testcase = document.createElement("testcase");
            	Attr testcaseName = document.createAttribute("name");
	            testcaseName.setValue(caseName);
	            testcase.setAttributeNode(testcaseName);

            	Attr caseFailure = document.createAttribute("failures");
	            caseFailure.setValue(Integer.toString(caseResult.getFirst()));
	            testcase.setAttributeNode(caseFailure);

	            Attr caseTotal = document.createAttribute("total");
	            caseTotal.setValue(Integer.toString(caseResult.getSecond()));
	            testcase.setAttributeNode(caseTotal);

	            testsuite.appendChild(testcase);

            }

        	root.appendChild(testsuite);
        }

        TransformerFactory transformerFactory = TransformerFactory.newInstance();
        Transformer transformer = transformerFactory.newTransformer();
        transformer.setOutputProperty(OutputKeys.INDENT, "yes");
        transformer.setOutputProperty("{http://xml.apache.org/xslt}indent-amount", "2");
        DOMSource domSource = new DOMSource(document);
        StreamResult streamResult = new StreamResult(new File(filePath));

        transformer.transform(domSource, streamResult);

        System.out.println("Done creating XML File");

	}
	private static void parseEachXML(HashMap<String, Pair<Pair<Integer, Integer>, HashMap<String, Pair<Integer, Integer>>>> flakey, File file) throws IOException, SAXException, ParserConfigurationException{

			DocumentBuilder dBuilder = DocumentBuilderFactory.newInstance()
			                             .newDocumentBuilder();

			Document doc = dBuilder.parse(file);
			identifyFailures(flakey, doc);
		
	}

	private static int convertMonth(String month) {
		if (month.equals("Jan")) {
			return 1;
		} else if (month.equals("Feb")) {
			return 2;
		} else if (month.equals("Mar")) {
			return 3;
		} else if (month.equals("Apr")) {
			return 4;
		} else if (month.equals("May")) {
			return 5;
		} else if (month.equals("Jun")) {
			return 6;
		} else if (month.equals("Jul")) {
			return 7;
		} else if (month.equals("Aug")) {
			return 8;
		} else if (month.equals("Sep")) {
			return 9;
		} else if (month.equals("Oct")) {
			return 10;
		} else if (month.equals("Nov")) {
			return 11;
		} else if (month.equals("Dec")) {
			return 12;
		}
		return 0;
	}

	private static boolean compareToPast(String date, int days) {
		int day = Integer.parseInt(date.substring(0, date.indexOf(" ")));
		date = date.substring(date.indexOf(" ") + 1);
		int month = convertMonth(date.substring(0, date.indexOf(" ")));
		int year = Integer.parseInt(date.substring(date.indexOf(" ") + 1));

		Calendar cal = Calendar.getInstance();
		cal.add(Calendar.DATE, -days);
		// String for date example: Tue May 14 14:22:48 PDT 2019
		String weekAgo = cal.getTime().toString();
		weekAgo = weekAgo.substring(weekAgo.indexOf(" ") + 1);
		int oldMonth = convertMonth(weekAgo.substring(0, weekAgo.indexOf(" ")));
		weekAgo = weekAgo.substring(weekAgo.indexOf(" ") + 1);
		int oldDay = Integer.parseInt(weekAgo.substring(0, weekAgo.indexOf(" ")));
		int oldYear = Integer.parseInt(weekAgo.substring(weekAgo.lastIndexOf(" ") + 1));

		if (year > oldYear || (year == oldYear && month > oldMonth) || (year == oldYear && month == oldMonth && day >= oldDay)){
			return true;
		}
		return false;

	}

	public static void testFlakey(String[] args) {
		try {
			//String command = "g=0; for n in $(gsutil ls gs://istio-circleci/master/*/*/artifacts/junit.xml); do foo=$(cut -d "," -f 2 <<< $(cut -d ":" -f 2 <<< $(gsutil stat $n | sed -n 3p))); gsutil cp $n " + '"' + "gs://istio-flakey-test/temp/out-$foo-$g.xml" + '"' + "; ((++g)); done";
			//Process process = Runtime.getRuntime().exec(command);
			HashMap<String, Pair<Pair<Integer, Integer>, HashMap<String, Pair<Integer, Integer>>>> flakey = new HashMap<>();
			//File dir = new File("temp");
			// File dir = new File(Paths.get(URI.create("gs://istio-flakey-test/temp")));
			// File[] foundFiles = dir.listFiles(new FilenameFilter() {
			//     public boolean accept(File dir, String name) {
			//         return name.startsWith("out-");
			//     }
			// });

			String outputFileName = "result.xml";
			int numDaysPast = 7;
			if (args.length >= 2) {
				outputFileName = args[0];
				numDaysPast = Integer.parseInt(args[1]);
			}

			Storage storage = StorageOptions.getDefaultInstance().getService();
			Page<Blob> blobs =
	    storage.list(
	        "istio-flakey-test", BlobListOption.currentDirectory(), BlobListOption.prefix("gs://istio-flakey-test/temp/out-"));
			for (Blob blob : blobs.iterateAll()) {
			  // do something with the blob
				String fileName = blob.getName();

				String fileContent = new String(blob.getContent());

				String date = fileName.substring(fileName.indexOf("-") + 1);
				date = date.substring(date.indexOf(" ") + 1);
				date = date.substring(0, date.lastIndexOf(" "));
				if (compareToPast(date, numDaysPast)) {
					DocumentBuilder dBuilder = DocumentBuilderFactory.newInstance()
			                             .newDocumentBuilder();
					InputSource is = new InputSource();
					is.setCharacterStream(new StringReader(fileContent));

					Document doc = dBuilder.parse(is);
					identifyFailures(flakey, doc);
					//parseEachXML(flakey, file);
				}
			}

			// for (File file : foundFiles) {
			// 	String fileName = file.getName();
			// 	String date = fileName.substring(fileName.indexOf("-") + 1);
			// 	date = date.substring(date.indexOf(" ") + 1);
			// 	date = date.substring(0, date.lastIndexOf(" "));
			// 	if (compareToPast(date, numDaysPast)) {
			// 		parseEachXML(flakey, file);
			// 	}
				
			// } 
			if (args.length >= 2) {
				outputFileName = args[0];
			}
			printFlakey(flakey, outputFileName);

		} catch (Exception e) {
			System.out.println(e.getMessage());
		}
	}

	public static String testHTML() {
		return "Test succeeded";
	}

	public static void main(String[] args) {
		// Storage storage = StorageOptions.getDefaultInstance().getService();
		// Page<Blob> blobs =
  //   storage.list(
  //       bucketName, BlobListOption.currentDirectory(), BlobListOption.prefix("gs://istio-flakey-test/temp/out-"));
		// for (Blob blob : blobs.iterateAll()) {
		//   // do something with the blob
		// 	String fileName = blob.getName();

		// 	String value = String(blob.getContent());
		// }
		testFlakey(args);
    }
}




